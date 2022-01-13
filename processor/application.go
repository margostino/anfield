package processor

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scrapper"
	"github.com/segmentio/kafka-go"
	"sync"
)

type Channels struct {
	commentary map[string]chan *domain.Commentary
	matchDate  map[string]chan string
	lineups    map[string]chan *domain.Lineups
}

type App struct {
	kafka         *kafka.Writer
	db            *db.Database
	scrapper      *scrapper.Scrapper
	configuration *configuration.Configuration
	channels      *Channels
	waitGroups    sync.Map
	logger        *Logger
}

func (a App) Start() error {
	var urls []string
	a.waitGroups = sync.Map{}
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
