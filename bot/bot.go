package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"log"
)

var config = common.GetConfig("./configuration/configuration.yml")

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
