package source

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"os"
)

var file *os.File

func Initialize() {
	config := context.Config()
	if context.ShouldUpdateData() {
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

func WriteOnFileIfUpdate(event *domain.Event) {
	for _, commentary := range event.Data {
		if file != nil {
			line := fmt.Sprintf("%s;%s;%s", event.Metadata.Date, commentary.Time, commentary.Comment)
			_, err := file.WriteString(line)
			common.Check(err)
		}
	}
}
