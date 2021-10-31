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
	browser := scrapper.Browser()
	defer browser.MustClose()
	normalizedUrls := make([]string, 0)
	baseUrl := config.Source.Url
	for _, url := range config.Realtime.Matches {
		normalizedUrls = append(normalizedUrls, baseUrl+url)
	}
	processor.Process(normalizedUrls)
}
