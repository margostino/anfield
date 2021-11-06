package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/processor"
	"log"
)

var config = context.BotConfig("./configuration/configuration.yml")
var bot *tgbotapi.BotAPI

func main() {
	context.Initialize()
	botApi, err := tgbotapi.NewBotAPI(config.Token)
	bot = botApi
	common.Check(err)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	poll(updates)
}

func poll(updates tgbotapi.UpdatesChannel) {
	welcome()
	go processor.Consume(config.Kafka.Topic, config.Kafka.Protocol, config.Kafka.Address)
	reply(updates)

}

func welcome() {
	for _, chatId := range context.Config().Bot.ChatIds {
		msg := tgbotapi.NewMessage(chatId, "Hi!!!")
		msg.ReplyMarkup = nil
		bot.Send(msg)
	}
}

func reply(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		replyMessage, replyMarkup := processor.Reply(&update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		msg.ReplyMarkup = replyMarkup
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
