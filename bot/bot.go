package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
	"strconv"
)

// TODO: implement DB in mem for history.
// TODO: set limit
var messagesHistory map[int64][]string

func NewBot(configuration *configuration.Configuration) *tgbotapi.BotAPI {
	bot, error := tgbotapi.NewBotAPI(configuration.Bot.Token)
	common.Check(error)
	//bot.Debug = true
	messagesHistory = make(map[int64][]string)
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func NewChannel() chan domain.User {
	return make(chan domain.User)
}

func (a App) poll(updates tgbotapi.UpdatesChannel) {
	a.Process(updates)
}

func (a App) listen() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := a.bot.GetUpdatesChan(updateConfig)
	return updates
}

func (a App) welcome() {
	for _, chatId := range a.configuration.Bot.ChatIds {
		msg := tgbotapi.NewMessage(chatId, "Hi!!!")
		msg.ReplyMarkup = nil
		a.bot.Send(msg)
	}
}

// Reply TODO: improve (reduce) amount ifs conditions. Make it generic
func (a App) Reply(update *tgbotapi.Update) (string, interface{}) {
	var markup interface{}
	var reply string
	message := update.Message.Text
	//username := update.Message.Chat.UserName
	userId := update.Message.Chat.ID

	if shouldStart(message) {
		user := getUserFrom(update)
		markup, reply = startReply(user)
		a.subscribe(user)
	}

	if shouldShowStats(message) {
		markup, reply = a.showStats(userId)
	}

	if isBuying(message) {
		// TODO: support asset+value in one command reply
		markup, reply = buyAssetQuestion()
	} else if shouldBuyAsset(getAllMessages(userId)) {
		markup, reply = buyUnitsQuestion()
	} else if shouldBuyAssetUnits(getAllMessages(userId)) {
		assetName := getLastMessage(userId)
		units, err := strconv.Atoi(message)
		if err != nil {
			markup, reply = buyInvalidUnits()
		} else {
			a.buy(userId, assetName, units)
		}
	}

	if reply == "" {
		markup, reply = echo(message) // TODO: tbd
	}

	appendPreviousMessage(userId, message) // TODO: improve this buffer strategy for back and forth

	return reply, markup
}

func (a App) Process(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go a.consume()
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		replyMessage, replyMarkup := a.Reply(&update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		msg.ReplyMarkup = replyMarkup
		//msg.ReplyToMessageID = update.Message.MessageID
		a.bot.Send(msg)
	}
}

func (a App) Send(message string) {
	// TODO: use subscription instead of static IDs from config
	//for _, chatId := range a.configuration.Bot.ChatIds {
	//	if IsFollowing(message, chatId) {
	//		msg := tgbotapi.NewMessage(chatId, message)
	//		msg.ReplyMarkup = nil
	//		a.bot.Send(msg) // TODO: filtering by subscription options
	//	}
	//}
}

// TODO: define fallback
func echo(message string) (interface{}, string) {
	return nil, message
}

func getUserFrom(update *tgbotapi.Update) *domain.User {
	user := update.Message.From
	return &domain.User{
		Id:        user.ID,
		Username:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func appendPreviousMessage(id int64, message string) {
	messagesHistory[id] = append(messagesHistory[id], message)
}

func getFirstMessage(id int64) string {
	if messages, ok := messagesHistory[id]; ok {
		if len(messages) > 0 {
			return messagesHistory[id][0]
		}
	}
	return ""
}

func getLastMessage(id int64) string {
	if messages, ok := messagesHistory[id]; ok {
		return messagesHistory[id][len(messages)-1]
	}
	return ""
}

func getAllMessages(id int64) []string {
	if messages, ok := messagesHistory[id]; ok {
		return messages
	}
	return nil
}
