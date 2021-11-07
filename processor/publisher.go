package processor

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func publish(event *Event) {

	b, err2 := json.Marshal(event)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(string(b))

	// TODO: send partial and not all.
	_, err := kafkaConnection.WriteMessages(
		kafka.Message{Value: []byte("new event: " + time.Now().String())},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
