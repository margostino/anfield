package main

import (
	"github.com/margostino/anfield/io"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/processor"
)

func main() {
	io.Initialize()
	processor.Initialize()
	webScrapper := processor.WebScrapper()

	file := io.File()
	if file != nil {
		defer file.Close()
	}
	defer webScrapper.Browser.MustClose()
	urls := processor.GetFinishedResults()
	processor.Process(urls)
	kafka.Close()
}
