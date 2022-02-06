// +build wireinject

package bot

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
)

func NewApp() (*App, error) {
	wire.Build(
		NewBot,
		NewActions,
		NewMessagesBuffer,
		db.NewDBConnection,
		configuration.GetConfig,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
