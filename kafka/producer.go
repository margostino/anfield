package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/margostino/anfield/domain"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

type Producer struct {
	Config *Config
	Client *kafka.Writer
}

func NewProducer(config *Config) *Producer {
	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	// make a writer that produces to topic-A, using the least-bytes distribution
	client := &kafka.Writer{
		Addr:     kafka.TCP(config.address),
		Topic:    config.topic,
		Balancer: &kafka.RoundRobin{},
	}
	return &Producer{
		Config: config,
		Client: client,
	}
}

func (w *Producer) Produce(metadata *domain.Metadata, commentary *domain.Commentary) {
	var message = domain.Message{
		Metadata: metadata,
		Data:     commentary,
	}
	messageBytes, _ := json.Marshal(message)
	id := strings.Split(metadata.Url, "/")[8]
	key := fmt.Sprintf("event-id-%s", id)
	err := w.Client.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func (w *Producer) Close() {
	if err := w.Client.Close(); err != nil {
		log.Fatal("failed to close kafka writer:", err)
	}
}
