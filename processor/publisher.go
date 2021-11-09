package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

func publish(metadata *Metadata, commentary *Commentary) {
	var message = Message{
		Metadata: metadata,
		Data:     commentary,
	}
	messageBytes, _ := json.Marshal(message)
	id := strings.Split(metadata.Url, "/")[8]
	key := fmt.Sprintf("event-id-%s", id)
	err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
