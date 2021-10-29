package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"strings"
)

func main() {
	var config = common.GetConfig("./configuration/configuration.yml")
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	resultsPage := browser.MustPage(config.Results.Url)
	results := resultsPage.MustElement(config.Results.Selector).MustElements(config.Results.Pattern)

	// TODO: parallelize event process
	for _, result := range results {
		eventUrl := result.MustProperty("href").String()
		infoPage := browser.MustPage(eventUrl + config.Matches.InfoUrlParam)
		startTimeDetail := infoPage.MustElement(config.Matches.InfoSelector).MustText()
		startTime := strings.Split(startTimeDetail, "\n")[0]
		summaryPage := browser.MustPage(eventUrl + config.Matches.CommentUrlParam)
		summaryPage.MustElement(config.Matches.MoreCommentsSelector).MustElement("*").MustClick()
		allComments := summaryPage.MustElement(config.Matches.CommentsSelector).MustText()
		events := strings.Split(allComments, "\n")
		eventName := strings.Split(eventUrl, "/")[7]
		fmt.Printf("======== START: %s ========\n", eventName)
		for _, event := range events {
			fmt.Printf("%s: %s\n", startTime, event)
		}
		fmt.Printf("======== END: %s ========\n", eventName)
	}

}
