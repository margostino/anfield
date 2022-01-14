package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var previousMessage string

func (a App) Reply(update *tgbotapi.Update) (string, interface{}) {
	var markup interface{}
	var reply string
	message := update.Message.Text
	//username := update.Message.Chat.UserName
	userId := update.Message.Chat.ID

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
		markup, reply = echo(message)
	}

	previousMessage = message

	return reply, markup
}

func (a App) Process(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		replyMessage, replyMarkup := a.Reply(&update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		msg.ReplyMarkup = replyMarkup
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func (a App) Send(message string) {
	// TODO: use subscription instead of static IDs from config
	for _, chatId := range a.configuration.Bot.ChatIds {
		if IsFollowing(message, chatId) {
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = nil
			Bot().Send(msg) // TODO: filtering by subscription options
		}
	}
}

// TODO: define fallback
func echo(message string) (interface{}, string) {
	return nil, message
}
