package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/margostino/anfield/domain"
)

//func Consume() {
//	for {
//		message, err := ReadMessage()
//
//		if err != nil {
//			break
//		}
//
//		commentary := concat(message)
//		scorer.CalculateScoring(message.Metadata.HomeTeam, message.Metadata.AwayTeam, message.Data)
//		bot.Send(commentary)
//		mongo.Insert(message)
//	}
//}

//func save(event *domain.Event) {
//	//mongo.Insert().Ins
//	eventLines := toString(event)
//	io.WriteOnFileIfUpdate(eventLines)
//}

//func toString(event *domain.Event) []string {
//	lines := make([]string, 0)
//	for _, commentary := range event.Data {
//		line := fmt.Sprintf("%s;%s;%s\n", event.Metadata.Date, commentary.Time, commentary.Comment)
//		lines = append(lines, line)
//	}
//	return lines
//}

func ReadMessage() (*domain.Message, error) {
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