package main

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/processor"
)

func main() {
	app, err := processor.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
	//kafka.NewWriter()
	//mongo.Initialize()
	//processor.Initialize()
	//urls := processor.GetUrlsResult()
	//processor.Process(urls)
	//kafka.Close()
	//mongo.Close()
	//processor.Close()
}
