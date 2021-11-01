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
	urls := make([]string, 0)
	if config.Realtime.Matches != nil {
		baseUrl := config.Source.Url
		for _, url := range config.Realtime.Matches {
			urls = append(urls, baseUrl+url)
		}
	} else {
		urls = scrapper.GetInProgressResults(browser)
	}

	processor.Process(urls)
}
