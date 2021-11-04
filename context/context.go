package context

import (
	"github.com/margostino/anfield/domain"
)

var configuration = GetConfig("./configuration/configuration.yml")

const (
	REALTIME = "realtime"
	BATCH    = "batch"
)

func Config() *domain.Configuration {
	return configuration
}

func ShouldUpdateData() bool {
	return configuration.Data.Update
}
