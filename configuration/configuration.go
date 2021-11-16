package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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

func ScoringRules() []Rule {
	return config.Rules
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

func NewDynamicRule(pattern string, pos int) *Rule {
	return &Rule{
		Pattern: pattern,
		Type:    DYNAMIC_RULE,
		Pos:     pos,
	}
}
