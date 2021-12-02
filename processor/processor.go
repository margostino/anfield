package processor

import (
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scrapper"
	"log"
	"strings"
	"sync"
)

var webScrapper *scrapper.Scrapper

//var waitGroups map[string]*sync.WaitGroup
var waitGroups = sync.Map{}
var stats = sync.Map{}
var metadataBuffer map[string]chan *domain.Metadata
var commentaryBuffer map[string]chan *domain.Commentary

func Initialize() {
	webScrapper = scrapper.New()
	//waitGroups = make(map[string]*sync.WaitGroup, 0)
	commentaryBuffer = make(map[string]chan *domain.Commentary)
	metadataBuffer = make(map[string]chan *domain.Metadata)
}

func WebScrapper() *scrapper.Scrapper {
	return webScrapper
}

func GetUrlsResult(mode string) []string {
	var urls []string

	if configuration.HasPredefinedEvents() {
		urls = getUrlsByConfig()
	} else {
		urls = getUrlsByScrapper(mode)
	}

	return urls
}

func Process(urls []string) {
	eventsToProcess := len(urls)

	if eventsToProcess == 0 {
		log.Println("URLs Not Found!")
	} else {
		wg := common.WaitGroup(len(urls))
		log.Println("Events to process", eventsToProcess)
		for _, url := range urls {
			waitGroups.Store(url, common.WaitGroup(3))
			//waitGroups[url] = common.WaitGroup(3)
			commentaryBuffer[url] = make(chan *domain.Commentary)
			metadataBuffer[url] = make(chan *domain.Metadata)
			go async(url, wg)
		}
		wg.Wait()
	}
}

func async(url string, waitGroup *sync.WaitGroup) {
	go produce(url)
	go consume(url)
	//waitGroups[url].Wait()
	wg, _ := waitGroups.Load(url)
	wg.(*sync.WaitGroup).Wait()
	waitGroup.Done()
}

func GetFinishedResults() []string {
	return GetUrlsResult(BATCH)
}

func GetInProgressResults() []string {
	return GetUrlsResult(REALTIME)
}

func getUrlsByConfig() []string {
	var urls []string
	matches := configuration.Events().Matches
	baseUrl := configuration.Scrapper().Url
	for _, url := range matches {
		urls = append(urls, baseUrl+url)
	}
	return urls
}

func getUrlsByScrapper(mode string) []string {
	var urls []string
	url := configuration.Scrapper().Url
	selector := configuration.Scrapper().MatchRowsSelector
	property := configuration.Scrapper().UrlProperty
	pattern := configuration.Scrapper().HrefPattern

	if mode == REALTIME {
		url += configuration.Scrapper().FixturesPath
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
			if true || mode == BATCH || (mode == REALTIME && inProgress(status)) {
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
