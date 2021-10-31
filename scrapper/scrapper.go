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

func GetUrlsResult(browser *rod.Browser, url string, selector string, pattern string) []string {
	var urls []string
	resultsPage := browser.MustPage(url)
	elements := resultsPage.MustElement(selector).MustElements(pattern)
	for _, element := range elements {
		urls = append(urls, element.MustProperty("href").String())
	}

	return urls

}

func GetEventDate(url string) string {
	selector := context.Config().Matches.InfoSelector
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
	moreCommentSelector := context.Config().Matches.MoreCommentsSelector
	commentSelector := context.Config().Matches.CommentsSelector
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
		isTime, _ := regexp.MatchString("([0-9]?'|[0-9]{2}'|[0-9]{2}\\+[0-9]+')$", event)

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
