package context

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func GetConfig(file string) *Configuration {
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
