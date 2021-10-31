package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/context"
	"log"
	"strings"
)

var config = context.GetConfig("./configuration/configuration.yml")

var subscriptions = make(map[string]string)
var matches = make([]string, 0)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	initializeMatches()

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
		replyMessage, replyMarkup := process(update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		msg.ReplyMarkup = replyMarkup
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func initializeMatches() {
	for _, match := range config.Realtime.Matches {
		matches = append(matches, strings.Split(match, "/")[1])
	}
}

func process(message string) (string, interface{}) {
	var markup interface{}
	var reply string
	if message == "/subscribe" {
		markup = getOptions()
		reply = "select a match to follow"
	} else {
		markup = nil
		reply = message
	}
	return reply, markup
}

func getOptions() interface{} {
	buttons := make([]tgbotapi.KeyboardButton, 0)
	for _, match := range matches {
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
