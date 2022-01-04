package main

import (
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/mongodb"
	"github.com/margostino/anfield/processor"
)

func main() {
	kafka.NewWriter()
	mongo.Initialize()
	processor.Initialize()
	urls := processor.GetUrlsResult()
	processor.Process(urls)
	kafka.Close()
	mongo.Close()
	processor.Close()
}
