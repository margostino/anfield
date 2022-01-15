// +build wireinject

package dataloader

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scorer"
)

func NewApp() (*App, error) {
	wire.Build(
		kafka.NewConfig,
		kafka.NewConsumer,
		db.NewDBConnection,
		scorer.NewScorer,
		configuration.GetConfig,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
