package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
)

type Channels struct {
	commentary map[string]chan *domain.Commentary
	matchDate  map[string]chan string
	lineups    map[string]chan *domain.Lineups
}

type App struct {
	kafka         *kafka.Consumer
	bot           *tgbotapi.BotAPI
	configuration *configuration.Configuration
}

func (a App) Start() error {
	a.welcome()
	//for _, match := range a.configuration.Realtime.Matches {
	//	matches = append(matches, strings.Split(match, "/")[1])
	//}
	updates := a.Listen()
	a.poll(updates)
	return nil // TODO tbd
}

func (a App) Close() {
	a.kafka.Close()
}
