package main

import (
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/processor"
	"github.com/margostino/anfield/scrapper"
	"github.com/margostino/anfield/source"
)

var config = context.GetConfig("./configuration/configuration.yml")

func main() {
	scrapper.Initialize()
	source.Initialize()
	file := source.File()
	if file != nil {
		defer file.Close()
	}
	browser := scrapper.Browser()
	defer browser.MustClose()
	url := config.Source.Url + config.Results.Url
	selector := config.Results.Selector
	pattern := config.Results.Pattern
	urls := scrapper.GetUrlsResult(browser, url, selector, pattern)
	processor.Process(urls)
}
