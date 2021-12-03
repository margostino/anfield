package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/bot"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/kafka"
)

func main() {
	bot.Initialize()
	kafka.NewReader(configuration.BotConsumerGroupId())
	updates := bot.Listen()
	poll(updates)
	kafka.Close()
}

func poll(updates tgbotapi.UpdatesChannel) {
	go bot.Consume()
	bot.Process(updates)
}
