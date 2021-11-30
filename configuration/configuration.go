package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var configFile = "./configuration.yml"
var rulesFile = "./rules.yml"
var config = getConfig(configFile, rulesFile)

func getConfig(configFile string, rulesFile string) *Configuration {
	var configuration Configuration
	var rules Rules
	unmarshal(configFile, &configuration)
	unmarshal(rulesFile, &rules)
	configuration.Rules = rules.ScoringRules
	return &configuration
}

func unmarshal(file string, out interface{}) {
	ymlFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, out)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func Scrapper() *ScrapperConfig {
	return config.Scrapper
}

func Bot() *BotConfig {
	return config.Bot
}

func BotConsumerGroupId() string {
	return config.Bot.KafkaConsumerGroupId
}

func DataLoaderConsumerGroupId() string {
	return config.DataLoader.KafkaConsumerGroupId
}

func Mongo() *MongoConfig {
	return config.Mongo
}

func Kafka() *KafkaConfig {
	return config.Kafka
}

func ChannelTimeout() time.Duration {
	return config.App.ChannelTimeout
}

func Realtime() *RealtimeConfig {
	return config.Realtime
}

func ScoringRules() []Rule {
	return config.Rules
}

func Source() *SourceConfig {
	return config.Source
}

func AppPath() string {
	return config.App.Path
}
