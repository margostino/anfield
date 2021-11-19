package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"log"
	"strings"
)

var bot *tgbotapi.BotAPI
var matches = make([]string, 0)
var subscriptions = make(map[int64][]string)
var following = make(map[int64][]string)

func Matches() []string {
	return matches
}

func Initialize() {
	bot = newBot()
	welcome()
	for _, match := range configuration.Realtime().Matches {
		matches = append(matches, strings.Split(match, "/")[1])
	}
}

func Listen() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := bot.GetUpdatesChan(updateConfig)
	return updates
}

func welcome() {
	for _, chatId := range configuration.Bot().ChatIds {
		msg := tgbotapi.NewMessage(chatId, "Hi!!!")
		msg.ReplyMarkup = nil
		bot.Send(msg)
	}
}

func Subscribe(userId int64, eventId string) {
	if !common.InSlice(eventId, subscriptions[userId]) {
		subscriptions[userId] = append(subscriptions[userId], eventId)
	}
}

func Follow(userId int64, player string) {
	following[userId] = append(following[userId], player)
}

func Unfollow(userId int64, player string) {
	following[userId] = common.Remove(following[userId], player)
}

func Bot() *tgbotapi.BotAPI {
	return bot
}

func newBot() *tgbotapi.BotAPI {
	bot, error := tgbotapi.NewBotAPI(configuration.Bot().Token)
	common.Check(error)
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func Following() map[int64][]string {
	return following
}
