package processor

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scrapper"
	"sync"
	"time"
)

type Channels struct {
	commentary map[string]chan *domain.Commentary
	matchDate  map[string]chan time.Time
	lineups    map[string]chan *domain.Lineups
}

type App struct {
	kafka         *kafka.Producer
	db            *db.Database
	scrapper      *scrapper.Scrapper
	configuration *configuration.Configuration
	channels      *Channels
	waitGroups    sync.Map
	logger        *Logger
}

func (a App) Start() error {
	var urls []string
	if a.configuration.HasPredefinedEvents() {
		urls = a.getUrlsByConfig()
	} else {
		urls = a.getUrlsByScrapper()
	}
	return a.Process(urls)
}

func (a App) Close() {
	a.db.Close()
	a.kafka.Close()
	a.scrapper.Close()
}
