package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/margostino/anfield/bot"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scorer"
)

func Consume() {
	for {
		message, err := readMessage()

		if err != nil {
			break
		}

		commentary := concat(message)
		scorer.CalculateScoring(message.Metadata.HomeTeam, message.Metadata.AwayTeam, message.Data)
		bot.Send(commentary)
	}
}

func readMessage() (*domain.Message, error) {
	var message domain.Message
	m, err := kafkaReader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	unmarshalError := json.Unmarshal(m.Value, &message)

	if unmarshalError != nil {
		fmt.Printf("Error when consuming message: %s\n", unmarshalError.Error())
	}

	//fmt.Printf("Message at offset %d: %s\n", m.Offset, string(m.Key))

	return &message, nil
}

func concat(message *domain.Message) string {
	return fmt.Sprintf("[%s] # %s\n", message.Data.Time, message.Data.Comment)
}
