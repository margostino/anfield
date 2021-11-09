package main

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/processor"
	"github.com/margostino/anfield/io"
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
		baseUrl := configuration.Source().Url
		for _, url := range matches {
			urls = append(urls, baseUrl+url)
		}
	} else {
		urls = processor.GetInProgressResults()
	}

	processor.Process(urls)
	processor.Close()
}
