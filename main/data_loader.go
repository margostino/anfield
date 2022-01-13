package main

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/dataloader"
)

func main() {
	app, err := dataloader.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
}
