// +build wireinject

package bot

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/db"
)

func NewApp() (*App, error) {
	wire.Build(
		db.NewDBConnection,
		NewBot,
		NewChannel,
		configuration.GetConfig,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
