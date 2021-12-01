package main

import (
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/processor"
)

func main() {
	kafka.NewWriter()
	processor.Initialize()
	webScrapper := processor.WebScrapper()
	defer webScrapper.Browser.MustClose()
	urls := processor.GetInProgressResults()
	processor.Process(urls)
	kafka.Close()
}
