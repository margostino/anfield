package source

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"os"
)

var file *os.File

func Initialize() {
	config := context.Config()
	if context.ShouldUpdateData()  {
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

func WriteOnFileIfUpdate(event string) {
	if file != nil {
		_, err := file.WriteString(event)
		common.Check(err)
	}
}