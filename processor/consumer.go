package processor

import (
	"fmt"
	"log"
)

func Consume() {

	batch := kafkaConnection.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	buffer := make([]byte, 10e3)                  // 10KB max per message

	for {
		noBytes, err := batch.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(string(buffer[:noBytes]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
}
