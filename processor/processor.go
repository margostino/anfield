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
var waitGroup *sync.WaitGroup
var stats = sync.Map{}
var matchDateBuffer map[string]chan string
var lineupsBuffer map[string]chan *domain.Lineups
var commentaryBuffer map[string]chan *domain.Commentary

func Initialize() {
	webScrapper = scrapper.New()
	//waitGroups = make(map[string]*sync.WaitGroup, 0)
	commentaryBuffer = make(map[string]chan *domain.Commentary)
	matchDateBuffer = make(map[string]chan string)
	lineupsBuffer = make(map[string]chan *domain.Lineups)
	InitializeLogger()
}

func Close() {
	webScrapper.Browser.MustClose()
}

func GetUrlsResult() []string {
	var urls []string

	if configuration.HasPredefinedEvents() {
		urls = getUrlsByConfig()
	} else {
		urls = getUrlsByScrapper()
	}

	return urls
}

func Process(urls []string) {
	eventsToProcess := len(urls)

	if eventsToProcess == 0 {
		log.Println("URLs Not Found!")
	} else {
		waitGroup = common.WaitGroup(len(urls) * (1 + 1 + 1 + 1))
		log.Println("Events to process: ", eventsToProcess)
		for _, url := range urls {
			initializeChannels(url)
			produce(url)
			go consume(url)
			//wait(url)
		}
		waitGroup.Wait()
	}
}

func wait(url string) {
	wg, _ := waitGroups.Load(url)
	wg.(*sync.WaitGroup).Wait()
}

func initializeChannels(url string) {
	waitGroups.Store(url, common.WaitGroup(4))
	commentaryBuffer[url] = make(chan *domain.Commentary)
	matchDateBuffer[url] = make(chan string)
	lineupsBuffer[url] = make(chan *domain.Lineups)
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

func getUrlsByScrapper() []string {
	baseUrl := configuration.Scrapper().Url
	fixturesUrl := baseUrl + configuration.Scrapper().FixturesPath // Matches in progress
	resultsUrl := baseUrl + configuration.Scrapper().ResultsPath   // Matches finished

	fixtureUrls := getMatchUrlsFrom(fixturesUrl)
	resultUrls := getMatchUrlsFrom(resultsUrl)

	return append(fixtureUrls, resultUrls...)
}

func getMatchUrlsFrom(url string) []string {
	var urls []string
	selector := configuration.Scrapper().MatchRowsSelector
	property := configuration.Scrapper().UrlProperty
	pattern := configuration.Scrapper().HrefPattern

	var prevSize = 0
	var tolerance = 35
	var currentSize = -1
	var equalsCounter = 0
	var elements rod.Elements

	webScrapper = webScrapper.GoPage(url)

	for {
		elements = webScrapper.DynamicElementsByPattern(selector, pattern)

		for _, element := range elements {
			//status := element.MustText() ---> inProgress(status) # For Realtime.
			url := element.MustProperty(property).String()
			if !common.InSlice(url, urls) {
				urls = append(urls, url)
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
