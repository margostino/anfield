// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package dataloader

import (
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scorer"
)

// Injectors from wire.go:

func NewApp() (*App, error) {
	configurationConfiguration := configuration.GetConfig()
	config := kafka.NewConfig(configurationConfiguration)
	consumer := kafka.NewConsumer(config)
	database := db.NewDBConnection()
	scorerScorer := scorer.NewScorer(configurationConfiguration)
	app := &App{
		kafka:         consumer,
		db:            database,
		scorer:        scorerScorer,
		configuration: configurationConfiguration,
	}
	return app, nil
}
