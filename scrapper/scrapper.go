package scrapper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"regexp"
	"strings"
	"time"
)

var browser *rod.Browser

func Initialize() {
	browser = rod.New().MustConnect()
}

func Browser() *rod.Browser {
	return browser
}

func GetFinishedResults(browser *rod.Browser) []string {
	return GetUrlsResult(browser, context.BATCH)
}

func GetInProgressResults(browser *rod.Browser) []string {
	return GetUrlsResult(browser, context.REALTIME)
}

func GetUrlsResult(browser *rod.Browser, mode string) []string {
	var urls []string
	var url, selector, pattern string

	if mode == context.REALTIME {
		selector = context.Config().Fixtures.Selector
		pattern = context.Config().Fixtures.Pattern
		url = context.Config().Source.Url + context.Config().Fixtures.Url
	} else {
		selector = context.Config().Results.Selector
		pattern = context.Config().Results.Pattern
		url = context.Config().Source.Url + context.Config().Results.Url
	}

	resultsPage := browser.MustPage(url)
	elements := resultsPage.MustElement(selector).MustElements(pattern)

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
	return isTimeCounter(prefix)
}

func GetEventDate(url string) string {
	selector := context.Config().Info.InfoSelector
	infoPage := browser.MustPage(url)
	startTimeDetail := infoPage.MustElement(selector).MustText()
	startTime := strings.Split(startTimeDetail, "\n")[0]
	day := strings.Split(startTime, " ")[0]
	month := strings.Split(startTime, " ")[1]
	year := strings.Split(startTime, " ")[2]
	normalizedStartTime := fmt.Sprintf("%s-%s-%s", year, common.NormalizeMonth(month), day)
	eventDate, _ := time.Parse("2006-01-02", normalizedStartTime)
	return eventDate.String()
}

// GetEvents TODO: read events as unbounded streams or until conditions (e.g. 90' time, message pattern, etc)
func GetEvents(date string, url string) *[]string {
	moreCommentSelector := context.Config().Commentary.MoreCommentsSelector
	commentSelector := context.Config().Commentary.CommentsSelector
	summaryPage := browser.MustPage(url)
	summaryPage.MustElement(moreCommentSelector).MustElement("*").MustClick()
	partialComments := summaryPage.MustElement(commentSelector).MustText()
	partialEvents := normalizeEvents(date, strings.Split(partialComments, "\n"))
	return &partialEvents
}

func normalizeEvents(eventDate string, events []string) []string {
	var time string
	var normalizedEvents = make([]string, 0)

	for _, event := range events {
		isTime := isTimeCounter(event)

		if isTime {
			time = event
		} else {
			normalizedEvent := fmt.Sprintf("%s;%s;%s\n", eventDate, time, event)
			normalizedEvents = append(normalizedEvents, normalizedEvent)
			time = ""
		}
	}
	return normalizedEvents
}

func isTimeCounter(value string) bool {
	isTime, _ := regexp.MatchString("([0-9]?'|[0-9]{2}'|[0-9]{2}\\+[0-9]+'|HT)$", value)
	return isTime
}
