package main

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/io"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/processor"
	"log"
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

	urls := make([]string, 0)
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
