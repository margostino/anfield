package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var configFile = "./configuration.yml"
var config = getConfig(configFile)

func getConfig(file string) *Configuration {
	var configuration Configuration
	ymlFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &configuration
}

func Config() *Configuration {
	return config
}

func Bot() *BotConfig {
	return config.Bot
}

func Kafka() *KafkaConfig {
	return config.Kafka
}

func Realtime() *RealtimeConfig {
	return config.Realtime
}

func Data() *DataConfig {
	return config.Data
}

func Fixture() *FixturesConfig {
	return config.Fixtures
}

func Results() *ResultsConfig {
	return config.Results
}

func Source() *SourceConfig {
	return config.Source
}

func Commentary() *CommentaryConfig {
	return config.Commentary
}

func Lineups() *LineupsConfig {
	return config.Lineups
}

func Info() *InfoConfig {
	return config.Info
}

func AppPath() string {
	return config.AppPath
}
