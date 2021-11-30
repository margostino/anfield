package bot

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	"github.com/margostino/anfield/scorer"
)

func Consume() {
	for {
		message, err := kafka.ReadMessage()

		if err != nil {
			break
		}

		commentary := concat(message)
		scorer.CalculateScoring(message.Metadata.HomeTeam, message.Metadata.AwayTeam, message.Data)
		Send(commentary)
	}
}

func concat(message *domain.Message) string {
	return fmt.Sprintf("[%s] # %s\n", message.Data.Time, message.Data.Comment)
}
