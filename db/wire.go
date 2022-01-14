// +build wireinject

package db

import (
	"github.com/google/wire"
	"github.com/margostino/anfield/configuration"
)

func NewDBConnection() *Database {
	panic(wire.Build(configuration.GetConfig, DefaultConnectionOpt, Connect))
}
