package main

import (
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/mongodb"
	"github.com/margostino/anfield/processor"
)

func main() {
	kafka.NewWriter()
	mongo.Initialize()
	processor.Initialize()
	webScrapper := processor.WebScrapper()
	defer webScrapper.Browser.MustClose()
	urls := processor.GetFinishedResults()
	processor.Process(urls)
	//urls = processor.GetFinishedResults()
	//processor.Process(urls)
	kafka.Close()
	mongo.Close()
}
