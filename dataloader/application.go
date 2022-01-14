package dataloader

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
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
	db            *db.Database
	configuration *configuration.Configuration
}

func (a App) Start() error {
	return a.Consume()
}

func (a App) Close() {
	a.db.Close()
	a.kafka.Close()
}
