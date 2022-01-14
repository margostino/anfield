package main

import (
	"github.com/margostino/anfield/bot"
	"github.com/margostino/anfield/common"
)

func main() {
	app, err := bot.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
}
