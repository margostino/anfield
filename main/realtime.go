package main

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/processor"
	"log"
)

func main() {
	kafka.NewWriter()
	processor.Initialize()
	webScrapper := processor.WebScrapper()
	urls := make([]string, 0)
	defer webScrapper.Browser.MustClose()
	matches := configuration.Realtime().Matches

	if matches != nil {
		baseUrl := configuration.Scrapper().Url
		for _, url := range matches {
			urls = append(urls, baseUrl+url)
		}
	} else {
		urls = processor.GetInProgressResults()
	}

	if len(urls) > 0 {
		processor.Process(urls)
	} else {
		log.Println("URLs Not Found!")
	}

	kafka.Close()
}
