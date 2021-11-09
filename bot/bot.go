package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/context"
)

func Reply(update *tgbotapi.Update) (string, interface{}) {
	var markup interface{}
	var reply string
	message := update.Message.Text
	username := update.Message.Chat.UserName
	if message == "/subscribe" {
		markup = getOptions()
		reply = "select a match to follow"
	} else {
		if common.InSlice(message, context.Matches()) {
			context.Subscribe(username, message)
			reply = "Done!"
		} else {
			reply = message
		}
		markup = nil

	}
	return reply, markup
}

func Send(message string) {
	for _, chatId := range configuration.Bot().ChatIds {
		msg := tgbotapi.NewMessage(chatId, message)
		msg.ReplyMarkup = nil
		context.Bot().Send(msg) // TODO: filtering by subscription options
	}
}

func getOptions() interface{} {
	buttons := make([]tgbotapi.KeyboardButton, 0)
	for _, match := range context.Matches() {
		button := tgbotapi.KeyboardButton{
			Text:            match,
			RequestContact:  false,
			RequestLocation: false,
		}
		buttons = append(buttons, button)
	}
	return &tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			buttons,
		},
	}
}
