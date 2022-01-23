package processor

import (
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"log"
	"strings"
	"sync"
	"time"
)

//var webScrapper *scrapper.Scrapper

//var waitGroups map[string]*sync.WaitGroup
//var waitGroups = sync.Map{}
var waitGroup *sync.WaitGroup
var stats = sync.Map{}

//var matchDateBuffer map[string]chan string
//var lineupsBuffer map[string]chan *domain.Lineups
//var commentaryBuffer map[string]chan *domain.Commentary

func Initialize() {
	//webScrapper = scrapper.New()
	//waitGroups = make(map[string]*sync.WaitGroup, 0)

	//commentaryBuffer = make(map[string]chan *domain.Commentary)
	//matchDateBuffer = make(map[string]chan string)
	//lineupsBuffer = make(map[string]chan *domain.Lineups)
}

func NewChannels() *Channels {
	return &Channels{
		commentary: make(map[string]chan *domain.Commentary),
		matchDate:  make(map[string]chan time.Time),
		lineups:    make(map[string]chan *domain.Lineups),
	}
}

func NewWaitGroups() sync.Map {
	return sync.Map{}
}

func (a App) InitializeChannels(url string) {
	a.waitGroups.Store(url, common.WaitGroup(4))
	a.channels.matchDate[url] = make(chan time.Time)
	a.channels.lineups[url] = make(chan *domain.Lineups)
	a.channels.commentary[url] = make(chan *domain.Commentary)
}

func (a App) Process(urls []string) error {
	eventsToProcess := len(urls)

	if eventsToProcess == 0 {
		log.Println("URLs Not Found!")
	} else {
		waitGroup = common.WaitGroup(len(urls) * (1 + 1 + 1 + 1))
		log.Println("Events to process: ", eventsToProcess)
		for _, url := range urls {
			a.InitializeChannels(url)
			a.produce(url)
			go a.consume(url)
			//wait(url)
		}
		waitGroup.Wait()
	}

	return nil // TODO: to be defined
}

func (a App) wait(url string) {
	wg, _ := a.waitGroups.Load(url)
	wg.(*sync.WaitGroup).Wait()
}

func (a App) getUrlsByConfig() []string {
	matches := a.configuration.Events.Matches
	baseUrl := a.configuration.Scrapper.Url
	var urls []string
	for _, url := range matches {
		urls = append(urls, baseUrl+url)
	}
	return urls
}

func (a App) getUrlsByScrapper() []string {
	baseUrl := a.configuration.Scrapper.Url
	fixturesUrl := baseUrl + a.configuration.Scrapper.FixturesPath // Matches in progress
	resultsUrl := baseUrl + a.configuration.Scrapper.ResultsPath   // Matches finished

	fixtureUrls := a.getMatchUrlsFrom(fixturesUrl)
	resultUrls := a.getMatchUrlsFrom(resultsUrl)

	return append(fixtureUrls, resultUrls...)
}

func (a App) getMatchUrlsFrom(url string) []string {
	selector := a.configuration.Scrapper.MatchRowsSelector
	property := a.configuration.Scrapper.UrlProperty
	pattern := a.configuration.Scrapper.HrefPattern

	var urls []string
	var prevSize = 0
	var tolerance = 35
	var currentSize = -1
	var equalsCounter = 0
	var elements rod.Elements

	webScrapper := a.scrapper.GoPage(url) // TODO: fix and reuse the scrapper in the struct instead of this

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
