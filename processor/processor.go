package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/scrapper"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"sync"
)

var webScrapper *scrapper.Scrapper
var waitGroups map[string]*sync.WaitGroup
var metadataBuffer map[string]chan *Metadata
var commentaryBuffer map[string]chan *Commentary
var kafkaConnection *kafka.Conn
var kafkaReader *kafka.Reader
var kafkaWiter *kafka.Writer

func Initialize() {
	webScrapper = scrapper.New()
	waitGroups = make(map[string]*sync.WaitGroup, 0)
	commentaryBuffer = make(map[string]chan *Commentary)
	metadataBuffer = make(map[string]chan *Metadata)
	kafkaWiter = NewKafkaWriter()
	kafkaReader = NewKafkaReader()
}

func Close() {
	if err := kafkaReader.Close(); err != nil {
		log.Fatal("failed to close kafka reader:", err)
	}
	if err := kafkaWiter.Close(); err != nil {
		log.Fatal("failed to close kafka writer:", err)
	}
}

func NewKafkaWriter() *kafka.Writer {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address

	// make a writer that produces to topic-A, using the least-bytes distribution
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return writer
}

func NewKafkaReader() *kafka.Reader {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address
	consumerGroupId := configuration.Kafka().ConsumerGroupId

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		GroupID:  consumerGroupId,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	return reader
}

func WebScrapper() *scrapper.Scrapper {
	return webScrapper
}

func Process(urls []string) {
	wg := common.WaitGroup(len(urls))
	for _, url := range urls {
		go async(url, wg)
	}
	wg.Wait()
}

func async(url string, waitGroup *sync.WaitGroup) {
	waitGroups[url] = common.WaitGroup(3)
	commentaryBuffer[url] = make(chan *Commentary)
	metadataBuffer[url] = make(chan *Metadata)

	go publishMetadata(url)
	go publishCommentary(url)
	// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
	go listen(url)

	waitGroups[url].Wait()
	waitGroup.Done()
}

func GetFinishedResults() []string {
	return GetUrlsResult(context.BATCH)
}

func GetInProgressResults() []string {
	return GetUrlsResult(context.REALTIME)
}

func GetUrlsResult(mode string) []string {
	var urls []string
	var url, selector, pattern string

	if mode == context.REALTIME {
		selector = configuration.Fixture().Selector
		pattern = configuration.Fixture().Pattern
		url = configuration.Source().Url + configuration.Fixture().Url
	} else {
		selector = configuration.Results().Selector
		pattern = configuration.Results().Pattern
		url = configuration.Source().Url + configuration.Results().Url
	}

	elements := webScrapper.GoPage(url).ElementsByPattern(selector, pattern)

	for _, element := range elements {
		status := element.MustText()
		if mode == context.BATCH || (mode == context.REALTIME && inProgress(status)) {
			urls = append(urls, element.MustProperty("href").String())
		}
	}

	return urls
}

func inProgress(status string) bool {
	prefix := strings.Split(status, "\n")[0]
	return common.IsTimeCounter(prefix)
}

func toString(event *Event) []string {
	lines := make([]string, 0)
	for _, commentary := range event.Data {
		line := fmt.Sprintf("%s;%s;%s\n", event.Metadata.Date, commentary.Time, commentary.Comment)
		lines = append(lines, line)
	}
	return lines
}
