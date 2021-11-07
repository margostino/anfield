package processor

import (
	goContext "context"
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/scrapper"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"sync"
	"time"
)

var webScrapper *scrapper.Scrapper
var waitGroups map[string]*sync.WaitGroup
var metadataBuffer map[string]chan *Metadata
var commentaryBuffer map[string]chan *Commentary
var kafkaConnection *kafka.Conn

func Initialize() {
	webScrapper = scrapper.New()
	waitGroups = make(map[string]*sync.WaitGroup, 0)
	commentaryBuffer = make(map[string]chan *Commentary)
	metadataBuffer = make(map[string]chan *Metadata)
	kafkaConnection = newKafkaConnection()
}

func Close() {
	if err := kafkaConnection.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func newKafkaConnection() *kafka.Conn {
	topic := configuration.Bot().Kafka.Topic
	protocol := configuration.Bot().Kafka.Protocol
	address := configuration.Bot().Kafka.Address
	partition := 0

	conn, err := kafka.DialLeader(goContext.Background(), protocol, address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	return conn
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
