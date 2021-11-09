package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/bot"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/processor"
)

func main() {
	bot.Initialize()
	processor.Initialize()
	updates := bot.Listen()
	poll(updates)
	kafka.Close()
}

func poll(updates tgbotapi.UpdatesChannel) {
	go kafka.Consume()
 	bot.Consume(updates)
}
