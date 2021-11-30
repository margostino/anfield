package processor

import (
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scrapper"
	"strings"
	"sync"
)

var webScrapper *scrapper.Scrapper
var waitGroups map[string]*sync.WaitGroup
var metadataBuffer map[string]chan *domain.Metadata
var commentaryBuffer map[string]chan *domain.Commentary

func Initialize() {
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
		waitGroups[url] = common.WaitGroup(3)
		go async(url, wg)
	}
	wg.Wait()
}

func async(url string, waitGroup *sync.WaitGroup) {
	commentaryBuffer[url] = make(chan *domain.Commentary)
	metadataBuffer[url] = make(chan *domain.Metadata)

	go produce(url)
	go consume(url)

	waitGroups[url].Wait()
	waitGroup.Done()
}

func GetFinishedResults() []string {
	return GetUrlsResult(BATCH)
}

func GetInProgressResults() []string {
	return GetUrlsResult(REALTIME)
}

func GetUrlsResult(mode string) []string {
	var urls []string
	selector := configuration.Scrapper().MatchRowsSelector
	pattern := configuration.Scrapper().HrefPattern
	property := configuration.Scrapper().UrlProperty
	url := configuration.Scrapper().Url

	if mode == REALTIME {
		url += configuration.Scrapper().FixturePath
	} else {
		url += configuration.Scrapper().ResultsPath
	}

	var prevSize = 0
	var tolerance = 35
	var currentSize = -1
	var equalsCounter = 0
	var elements rod.Elements

	webScrapper = webScrapper.GoPage(url)

	for {
		elements = webScrapper.DynamicElementsByPattern(selector, pattern)

		for _, element := range elements {
			status := element.MustText()
			if mode == BATCH || (mode == REALTIME && inProgress(status)) {
				url := element.MustProperty(property).String()
				if !common.InSlice(url, urls) {
					urls = append(urls, url)
				}
			}
		}

		prevSize, currentSize = currentSize, len(urls)

		if prevSize == currentSize {
			equalsCounter += 1
		} else {
			equalsCounter = 0
		}

		if equalsCounter == tolerance {
			break
		}
	}

	return urls
}

func inProgress(status string) bool {
	prefix := strings.Split(status, "\n")[0]
	return common.IsTimeCounter(prefix)
}
