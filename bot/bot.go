package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"log"
)

var previousMessage string

func Reply(update *tgbotapi.Update) (string, interface{}) {
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
		markup, reply = followReply()
	} else if shouldFollowPlayer(previousMessage) {
		markup, reply = playerFollowerReply(message, userId)
	}

	if shouldUnfollow(message) {
		markup, reply = unfollowReply()
	} else if shouldUnfollowPlayer(previousMessage) {
		markup, reply = playerUnfollowerReply(message, userId)
	}

	if reply == "" {
		markup, reply = echo(message)
	}

	previousMessage = message

	return reply, markup
}

func Consume(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		replyMessage, replyMarkup := Reply(&update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		msg.ReplyMarkup = replyMarkup
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func Send(message string) {
	// TODO: use subscription instead of static IDs from config
	for _, chatId := range configuration.Bot().ChatIds {
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
