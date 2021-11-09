package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/processor"
	"log"
)

var bot *tgbotapi.BotAPI

func main() {
	context.Initialize()
	processor.Initialize()
	bot = context.Bot()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := bot.GetUpdatesChan(updateConfig)
	welcome()
	poll(updates)
	processor.Close()
}

func poll(updates tgbotapi.UpdatesChannel) {
	go processor.Consume()
	reply(updates)
}

func welcome() {
	for _, chatId := range configuration.Bot().ChatIds {
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
