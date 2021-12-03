package bot

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
)

func Consume() {
	for {
		message, err := kafka.ReadMessage()

		if err != nil {
			break
		}

		commentary := concat(message)
		Send(commentary)
	}
}

func concat(message *domain.Message) string {
	return fmt.Sprintf("[%s] # %s\n", message.Data.Time, message.Data.Comment)
}
