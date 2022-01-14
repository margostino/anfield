package main

import (
	"github.com/margostino/anfield/bot"
	"github.com/margostino/anfield/common"
)

//	Application to send updates to subscribers: new scoring, data change, lineups, team news, etc.
// 	Every new should be consumed from different topics.
// 	Every user can be subscribed to zero or more topics.
// 	Every user can trade assets during a game in realtime.
func main() {
	app, err := bot.NewApp()
	common.Check(err)
	app.Start()
	app.Close()
}
