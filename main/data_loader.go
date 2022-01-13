package main

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/dataloader"
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/db"
)

func main() {
	consumerGroupId := configuration.DataLoaderConsumerGroupId()
	kafka.NewReader(consumerGroupId)
	mongo.Initialize()
	dataloader.Consume()
	kafka.Close()
	mongo.Close()
}