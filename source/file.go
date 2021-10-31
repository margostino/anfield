package source

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"os"
)

var file *os.File

func Initialize() {
	config := context.Config()
	if config.Data.Update {
		filename := config.AppPath + config.Data.Matches
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

func WriteOnFile(data string) {
	if context.ShouldUpdateData() {
		_, err := file.WriteString(data)
		common.Check(err)
	}
}
