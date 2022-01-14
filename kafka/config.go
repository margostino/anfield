package kafka

import (
	"github.com/margostino/anfield/configuration"
)

type Config struct {
	topic           string
	address         string
	consumerGroupId string
}

func NewConfig(configuration *configuration.Configuration) *Config {
	return &Config{
		address:         configuration.Kafka.Address,
		topic:           configuration.Kafka.Topic,
		consumerGroupId: configuration.Kafka.ConsumerGroupId,
	}
}
