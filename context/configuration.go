package context

import (
	"github.com/margostino/anfield/domain"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func GetConfig(file string) *domain.Configuration {
	var configuration domain.Configuration
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

func BotConfig(file string) *domain.BotConfig {
	return GetConfig(file).Bot
}
