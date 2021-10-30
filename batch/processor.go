package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/anfield/common"
	"os"
	"regexp"
	"strings"
	"time"
)

var config = common.GetConfig("./configuration/configuration.yml")

func main() {
	file := openFile()
	if file != nil {
		defer file.Close()
	}
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
	if config.Data.Update {
		filePath := config.AppPath + config.Data.Matches
		file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
		common.Check(err)
		return file
	}
	return nil
}

func openBrowser() *rod.Browser {
	browser := rod.New().MustConnect()
	return browser
}

func getResults(browser *rod.Browser) rod.Elements {
	resultsPage := browser.MustPage(config.Source.Url + config.Results.Url)
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
	var time string
	for _, event := range *events {
		isTime, _ := regexp.MatchString("([0-9]?'|[0-9]{2}'|[0-9]{2}\\+[0-9]+')$", event)

		if isTime {
			time = event
		} else {
			line := fmt.Sprintf("%s;%s;%s\n", eventDate, time, event)
			fmt.Println(line)
			persist(file, line)
			time = ""
		}
	}
}

func persist(file *os.File, data string) {
	if config.Data.Update {
		_, err := file.WriteString(data)
		common.Check(err)
	}
}
