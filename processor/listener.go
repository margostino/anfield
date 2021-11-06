package processor

import (
	"context"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/source"
	"github.com/segmentio/kafka-go"
	"time"
)

// TODO: calculate stats, bot sender
// This aggregation in consumer should happen once by URL/Event
func listen(url string) {

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
			publish(event.Metadata, commentary)
		}

	}
}
