package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"os"
	"strings"
	"time"
)

var appConfig = common.GetConfig("./configuration/configuration.yml")
var config = appConfig.Batch

func main() {
	file := openFile()
	defer file.Close()
	browser := openBrowser()
	defer browser.MustClose()
	results := getResults(browser)

	// TODO: parallelize event process
	for _, result := range results {
		eventUrl := result.MustProperty("href").String()
		eventDate := getEventDate(browser, eventUrl)
		eventName, events := getEvents(browser, eventUrl)
		fmt.Printf("======== START: %s ========\n", eventName)
		save(eventDate, events, file)
		fmt.Printf("======== END: %s ========\n", eventName)
	}

}

func openFile() *os.File {
	filePath := appConfig.AppPath + appConfig.Data.Matches
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	common.Check(err)
	return file
}

func openBrowser() *rod.Browser {
	browser := rod.New().MustConnect()
	return browser
}

func getResults(browser *rod.Browser) rod.Elements {
	resultsPage := browser.MustPage(config.Results.Url)
	return resultsPage.MustElement(config.Results.Selector).MustElements(config.Results.Pattern)
}

func getEventDate(browser *rod.Browser, eventUrl string) string {
	infoPage := browser.MustPage(eventUrl + config.Matches.InfoUrlParam)
	startTimeDetail := infoPage.MustElement(config.Matches.InfoSelector).MustText()
	startTime := strings.Split(startTimeDetail, "\n")[0]
	day := strings.Split(startTime, " ")[0]
	month := strings.Split(startTime, " ")[1]
	year := strings.Split(startTime, " ")[2]
	normalizedStartTime := fmt.Sprintf("%s-%s-%s", year, common.NormalizeMonth(month), day)
	eventDate, _ := time.Parse("2006-01-02", normalizedStartTime)
	return eventDate.String()
}

func getEvents(browser *rod.Browser, eventUrl string) (string, *[]string) {
	summaryPage := browser.MustPage(eventUrl + config.Matches.CommentUrlParam)
	summaryPage.MustElement(config.Matches.MoreCommentsSelector).MustElement("*").MustClick()
	allComments := summaryPage.MustElement(config.Matches.CommentsSelector).MustText()
	events := strings.Split(allComments, "\n")
	eventName := strings.Split(eventUrl, "/")[7]
	return eventName, &events
}

func save(eventDate string, events *[]string, file *os.File) {
	for _, event := range *events {
		line := fmt.Sprintf("%s;%s\n", eventDate, event)
		fmt.Println(line)
		_, err := file.WriteString(line)
		common.Check(err)
	}
}
