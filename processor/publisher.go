package processor

import (
	"context"
	context2 "github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func publish(metadata *domain.Metadata, commentary *domain.Commentary) {
	// to produce messages
	topic := context2.Config().Bot.Kafka.Topic
	address := context2.Config().Bot.Kafka.Address
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("new event: " + time.Now().String())},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
