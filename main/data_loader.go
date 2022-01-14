package main

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/dataloader"
)

//	Application to consume every event (metadata + commentary) in topic and aggregate and calculate scoring.
//	Every aggregation and scoring is sent to another topic which is consumed by Bot Application.
func main() {
	app, err := dataloader.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
}
