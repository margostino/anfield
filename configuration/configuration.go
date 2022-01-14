package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func GetConfig() *Configuration {
	var configuration Configuration
	var rules Rules
	unmarshal("./configuration.yml", &configuration)
	unmarshal("./rules.yml", &rules)
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

func (c *Configuration) BotConsumerGroupId() string {
	return c.Bot.KafkaConsumerGroupId
}

func (c *Configuration) DataLoaderConsumerGroupId() string {
	return c.DataLoader.KafkaConsumerGroupId
}

func (c *Configuration) ChannelTimeout() time.Duration {
	return c.App.ChannelTimeout
}

func (c *Configuration) HasPredefinedEvents() bool {
	return c.Events.Matches != nil
}

func (c *Configuration) AppPath() string {
	return c.App.Path
}
