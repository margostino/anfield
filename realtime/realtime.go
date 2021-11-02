package main

import (
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/processor"
	"github.com/margostino/anfield/scrapper"
	"github.com/margostino/anfield/source"
)

var config = context.GetConfig("./configuration/configuration.yml")

// TODO: spawn multiple event consumer in parallel
func main() {
	scrapper.Initialize()
	source.Initialize()
	browser := scrapper.Browser()
	defer browser.MustClose()

	//summaryPage := browser.MustPage("https://www.livescore.com/en/football/england/premier-league/wolverhampton-wanderers-vs-everton/450665/?tab=summary-commentary")
	//summaryPage.MustElement(".Summary_toggleAllComments__1M7FV").MustElement("*").MustClick()
	//rawEvents := summaryPage.MustElement(".Summary_blockWrapper__1P4fu").MustText()
	////rawEvents := summaryPage.MustElement(".Summary_blockWrapper__1P4fu").MustText()
	////rawEvents := summaryPage.MustElements("*")
	//events := strings.Split("nil", "\n")
	//println(rawEvents)
	////for _, element := range rawEvents {
	////	println(element.MustText())
	////}
	//println(events)

	urls := make([]string, 0)
	if config.Realtime.Matches != nil {
		baseUrl := config.Source.Url
		for _, url := range config.Realtime.Matches {
			urls = append(urls, baseUrl+url)
		}
	} else {
		urls = scrapper.GetInProgressResults()
	}
	processor.Process(urls)
}
