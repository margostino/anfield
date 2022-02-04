// +build wireinject

package bot

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
)

func NewApp() (*App, error) {
	wire.Build(
		NewBot,
		NewChannel,
		NewActions,
		db.NewDBConnection,
		configuration.GetConfig,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
