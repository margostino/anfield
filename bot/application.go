package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
)

type App struct {
	bot           *tgbotapi.BotAPI
	db            *db.Database
	configuration *configuration.Configuration
	subscriptions chan domain.User
}

func (a App) Start() error {
	//a.welcome()
	//for _, match := range a.configuration.Realtime.Matches {
	//	matches = append(matches, strings.Split(match, "/")[1])
	//}
	updates := a.listen()
	a.poll(updates)
	return nil // TODO tbd
}

func (a App) Close() {
	a.db.Close()
}
