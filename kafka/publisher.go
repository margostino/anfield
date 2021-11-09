package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/margostino/anfield/processor"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

func Publish(metadata *processor.Metadata, commentary *processor.Commentary) {
	var message = processor.Message{
		Metadata: metadata,
		Data:     commentary,
	}
	messageBytes, _ := json.Marshal(message)
	id := strings.Split(metadata.Url, "/")[8]
	key := fmt.Sprintf("event-id-%s", id)
	err := processor.KafkaWriter().WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
