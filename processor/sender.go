package processor

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
)

func Send(message string) {
	for _, chatId := range configuration.Bot().ChatIds {
		msg := tgbotapi.NewMessage(chatId, message)
		msg.ReplyMarkup = nil
		context.Bot().Send(msg) // TODO: filtering by subscription options
	}
}
