package main

import (
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/mongodb"
	"github.com/margostino/anfield/processor"
)

// TODO: unify with batch and apply strategy on results in progress
func main() {
	kafka.NewWriter()
	mongo.Initialize()
	processor.Initialize()
	urls := processor.GetInProgressResults()
	processor.Process(urls)
	kafka.Close()
	mongo.Close()
	processor.Close()
}
