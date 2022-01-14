package bot

import (
	"fmt"
	"github.com/margostino/anfield/domain"
)

func (a App) Consume() {
	for {
		message, err := a.kafka.ReadMessage()

		if err != nil {
			break
		}

		commentary := concat(message)
		a.Send(commentary)
	}
}

func concat(message *domain.Message) string {
	return fmt.Sprintf("[%s] # %s\n", message.Data.Time, message.Data.Comment)
}
