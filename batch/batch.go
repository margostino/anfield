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
	urls := scrapper.GetFinishedResults(browser)
	processor.Process(urls)
}
