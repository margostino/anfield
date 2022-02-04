package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
)

// TODO: implement DB in mem for history.
// TODO: set limit
var messagesHistory map[int][]string

func NewBot(configuration *configuration.Configuration) *tgbotapi.BotAPI {
	bot, error := tgbotapi.NewBotAPI(configuration.Bot.Token)
	common.Check(error)
	//bot.Debug = true
	messagesHistory = make(map[int][]string)
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
// Reply: TODO: support one shot command values
func (a App) Reply(update *tgbotapi.Update) (string, interface{}) {
	var reply string
	var markup interface{}
	message := update.Message.Text
	userId := update.Message.From.ID

	if isCommand(message) {
		cleanupAllPreviousMessagesBy(userId)
	}

	if shouldStart(message) {
		user := getUserFrom(update)
		markup, reply = replyStart(user)
		a.subscribe(user)
	}

	if shouldShowStats(message) {
		markup, reply = replyStats(userId)
	}

	if isBuying(message) {
		// TODO: support asset+value in one command reply
		markup, reply = buyAssetValueInstruction()
	} else if shouldBuyAsset(getAllMessages(userId)) {
		assetName, units, err := extractTransactionFrom(message)
		if err != nil {
			markup, reply = nil, err.Error()
		} else {
			a.buy(userId, assetName, units)
			cleanupAllPreviousMessagesBy(userId)
		}
	}

	if reply == "" {
		markup, reply = fallback()
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

func fallback() (interface{}, string) {
	return nil, "Input is not expected"
}

func getUserFrom(update *tgbotapi.Update) *domain.User {
	from := update.Message.From
	return &domain.User{
		SocialId:  from.ID,
		Username:  from.UserName,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	}
}

func cleanupAllPreviousMessagesBy(userId int) {
	messagesHistory[userId] = make([]string, 0)
}

func appendPreviousMessage(id int, message string) {
	messagesHistory[id] = append(messagesHistory[id], message)
}

func isCommand(message string) bool {
	return message[0:1] == "/"
}

func getFirstMessage(id int) string {
	if messages, ok := messagesHistory[id]; ok {
		if len(messages) > 0 {
			return messagesHistory[id][0]
		}
	}
	return ""
}

func getLastMessage(id int) string {
	if messages, ok := messagesHistory[id]; ok {
		return messagesHistory[id][len(messages)-1]
	}
	return ""
}

func getAllMessages(id int) []string {
	if messages, ok := messagesHistory[id]; ok {
		return messages
	}
	return nil
}
