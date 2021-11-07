package main

import (
	"github.com/margostino/anfield/processor"
	"github.com/margostino/anfield/source"
)

func main() {
	source.Initialize()
	processor.Initialize()
	webScrapper := processor.WebScrapper()

	file := source.File()
	if file != nil {
		defer file.Close()
	}
	defer webScrapper.Browser.MustClose()
	urls := processor.GetFinishedResults()
	processor.Process(urls)
	processor.Close()
}
