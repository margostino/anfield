package main

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/processor"
)

//	Application to process all events (realtime and batch mode).
//	Every game is processed in a different goroutines:
//	=> Producer for lineups, commentary and info.
//	=> Consumer and Merged as sink fashion processor which sends to Kafka topic the event (metadata + commentary)
func main() {
	app, err := processor.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
}
