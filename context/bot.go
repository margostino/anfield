package context

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"log"
	"strings"
)

var bot *tgbotapi.BotAPI
var matches = make([]string, 0)
var subscriptions = make(map[string]string)

func Subscriptions() map[string]string {
	return subscriptions
}

func Matches() []string {
	return matches
}

func Initialize() {
	newBot()
	for _, match := range configuration.Realtime().Matches {
		matches = append(matches, strings.Split(match, "/")[1])
	}
}

func Subscribe(username string, eventId string) {
	subscriptions[username] = eventId
}

func Bot() *tgbotapi.BotAPI {
	return bot
}

func newBot() {
	b, error := tgbotapi.NewBotAPI(configuration.Bot().Token)
	bot = b
	common.Check(error)
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}
