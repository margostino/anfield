package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	Config *Config
	Client *kafka.Reader
}

func NewConsumer(config *Config) *Consumer {
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	client := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.address},
		GroupID: config.consumerGroupId,
		Topic:   config.topic,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
	})
	return &Consumer{
		Config: config,
		Client: client,
	}
}

func (r *Consumer) Consume(event interface{}) error {
	m, err := r.Client.ReadMessage(context.Background())
	if err != nil {
		return err
	}

	unmarshalError := json.Unmarshal(m.Value, &event)

	if unmarshalError != nil {
		fmt.Printf("Error when consuming message: %s\n", unmarshalError.Error())
	}

	//fmt.Printf("Message at offset %d: %s\n", m.Offset, string(m.Key))

	return nil
}

func (r *Consumer) Close() {
	if err := r.Client.Close(); err != nil {
		log.Fatal("failed to close kafka reader bus:", err)
	}
}
