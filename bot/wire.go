// +build wireinject

package bot

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/kafka"
)

func NewApp() (*App, error) {
	wire.Build(
		kafka.NewConfig,
		kafka.NewConsumer,
		NewBot,
		configuration.GetConfig,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
