package processor

import (
	"context"
	"fmt"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/source"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// TODO: calculate stats, bot sender
// This aggregation in consumer should happen once by URL/Event
func consume(url string) {

	// to create topics when auto.create.topics.enable='true'
	_, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "my-topic", 0)
	if err != nil {
		panic(err.Error())
	}

	metadata := <-metadataBuffer[url]
	event := NewEvent(metadata)
	commentaryLoop(event)
	source.WriteOnFileIfUpdate(event)
	waitGroups[url].Done()
}

func NewEvent(metadata *domain.Metadata) *domain.Event {
	return &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
}

func commentaryLoop(event *domain.Event) {
	url := event.Metadata.Url
	h2h := event.Metadata.H2H
	for {
		commentary := <-commentaryBuffer[url]
		event.Data = append(event.Data, commentary)

		time.Sleep(100 * time.Millisecond) // TODO: configurable

		if end(commentary) {
			break
		} else {
			printCommentary(h2h, commentary)
			publishEvent(event.Metadata, commentary)
		}

	}
}

// TODO: wip
func publishEvent(metadata *domain.Metadata, commentary *domain.Commentary)  {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func printCommentary(h2h string, commentary *domain.Commentary) {
	if isTimedComment(commentary) {
		fmt.Printf("[%s] # %s\n", h2h, commentary.Comment)
	} else {
		fmt.Printf("[%s] # %s - %s\n", h2h, commentary.Time, commentary.Comment)
	}
}

func end(commentary *domain.Commentary) bool {
	if commentary == nil || (commentary.Time == "end" && commentary.Comment == "end") {
		return true
	}
	return false
}

func isTimedComment(commentary *domain.Commentary) bool {
	if commentary != nil {
		if commentary.Time == "" && commentary.Comment != "" {
			return true
		}
	}
	return false
}
