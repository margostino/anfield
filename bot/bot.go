package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
)

var matches = make([]string, 0)
var subscriptions = make(map[int64][]string)
var following = make(map[int64][]string)
var previousMessage string

func Matches() []string {
	return matches
}

func NewBot(configuration *configuration.Configuration) *tgbotapi.BotAPI {
	bot, error := tgbotapi.NewBotAPI(configuration.Bot.Token)
	common.Check(error)
	//bot.Debug = true
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

	if isSubscription(message) {
		markup, reply = subscriptionReply()
	} else if shouldSubscribeToMatch(previousMessage, message) {
		markup, reply = matchSubscriptionReply(message, userId)
	}

	if shouldFollow(message) {
		markup, reply = followQuestion()
	} else if shouldFollowPlayer(previousMessage) {
		markup, reply = followReply(message, userId)
	}

	if shouldUnfollow(message) {
		markup, reply = unfollowQuestion()
	} else if shouldUnfollowPlayer(previousMessage) {
		markup, reply = unfollowReply(message, userId)
	}

	if shouldShowStats(message) {
		markup, reply = a.showStats(userId)
	}

	if reply == "" {
		markup, reply = echo(message) // TODO: tbd
	}

	previousMessage = message

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
	for _, chatId := range a.configuration.Bot.ChatIds {
		if IsFollowing(message, chatId) {
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = nil
			a.bot.Send(msg) // TODO: filtering by subscription options
		}
	}
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
