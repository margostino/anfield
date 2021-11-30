package main

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/mongodb"
	"github.com/margostino/anfield/processor"
)

func main() {
	urls := make([]string, 0)
	kafka.NewWriter()
	mongo.Initialize()
	processor.Initialize()
	webScrapper := processor.WebScrapper()
	defer webScrapper.Browser.MustClose()

	matches := configuration.Realtime().Matches
	if matches != nil {
		baseUrl := configuration.Scrapper().Url
		for _, url := range matches {
			urls = append(urls, baseUrl+url)
		}
	} else {
		urls = processor.GetFinishedResults()
	}


	processor.Process(urls)
	kafka.Close()
	mongo.Close()
}
