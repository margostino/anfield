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
}
