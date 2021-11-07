package processor

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func publish(event *Event) {
	eventJson, _ := json.Marshal(event)
	// TODO: send partial and not all.
	err := kafkaWiter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("event-id-" + string(len(event.Data))),
			Value: []byte(string(eventJson)),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
