package io

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"os"
)

var file *os.File

func Initialize() {
	if configuration.Source().Update {
		filename := configuration.AppPath() + configuration.Source().MatchesPath
		f, err := os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
		common.Check(err)
		file = f
	} else {
		file = nil
	}
}

func File() *os.File {
	return file
}

func WriteOnFileIfUpdate(lines []string) {
	for _, line := range lines {
		if file != nil {
			_, err := file.WriteString(line)
			common.Check(err)
		}
	}
}
