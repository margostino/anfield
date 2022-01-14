// +build wireinject

package processor

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scrapper"
)

func NewApp() (*App, error) {
	wire.Build(
		kafka.NewConfig,
		kafka.NewProducer,
		db.NewDBConnection,
		scrapper.New,
		configuration.GetConfig,
		NewChannels,
		NewWaitGroups,
		NewLogger,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
