package processor

import (
	"context"
	"encoding/json"
	"fmt"
)

func Consume() {
	for {
		var message Message
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		unmarshalError := json.Unmarshal(m.Value, &message)

		if unmarshalError != nil {
			fmt.Printf("Error when consuming message: %s\n", unmarshalError.Error())
		}

		commentary := fmt.Sprintf("[%s] # %s\n", message.Data.Time, message.Data.Comment)
		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), commentary)
	}
}
