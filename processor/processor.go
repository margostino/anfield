package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scrapper"
	"strings"
	"sync"
)

var webScrapper *scrapper.Scrapper
var waitGroups map[string]*sync.WaitGroup
var metadataBuffer map[string]chan *domain.Metadata
var commentaryBuffer map[string]chan *domain.Commentary

func Initialize() {
	kafka.Initialize()
	webScrapper = scrapper.New()
	waitGroups = make(map[string]*sync.WaitGroup, 0)
	commentaryBuffer = make(map[string]chan *domain.Commentary)
	metadataBuffer = make(map[string]chan *domain.Metadata)
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
	commentaryBuffer[url] = make(chan *domain.Commentary)
	metadataBuffer[url] = make(chan *domain.Metadata)

	go produce(url)
	go consume(url)

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

func toString(event *domain.Event) []string {
	lines := make([]string, 0)
	for _, commentary := range event.Data {
		line := fmt.Sprintf("%s;%s;%s\n", event.Metadata.Date, commentary.Time, commentary.Comment)
		lines = append(lines, line)
	}
	return lines
}
