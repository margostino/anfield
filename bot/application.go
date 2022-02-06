package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
)

type App struct {
	actions       []Action
	db            *db.Database
	bot           *tgbotapi.BotAPI
	configuration *configuration.Configuration
	messageBuffer map[int]string // TODO: currently only Bot support 1 level (ask/reply). Evaluate increase the level.
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
