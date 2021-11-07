package processor

import (
	"context"
	"fmt"
)

func Consume() {
	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
