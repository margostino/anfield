package scrapper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
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
func GetEvents(url string) (string, *[]string) {
	moreCommentSelector := context.Config().Matches.MoreCommentsSelector
	commentSelector := context.Config().Matches.CommentsSelector
	summaryPage := browser.MustPage(url)
	summaryPage.MustElement(moreCommentSelector).MustElement("*").MustClick()
	allComments := summaryPage.MustElement(commentSelector).MustText()
	events := strings.Split(allComments, "\n")
	eventName := strings.Split(url, "/")[7]
	return eventName, &events
}
