package processor

import (
	"github.com/margostino/anfield/io"
	"github.com/margostino/anfield/kafka"
	"time"
)

// TODO: calculate stats, bot sender
// This aggregation in consumer should happen once by URL/Event
func listen(url string) {
	metadata := <-metadataBuffer[url]
	event := NewEvent(metadata)
	commentaryLoop(event)
	eventLines := toString(event)
	io.WriteOnFileIfUpdate(eventLines)
	waitGroups[url].Done()
}

func NewEvent(metadata *Metadata) *Event {
	return &Event{
		Metadata: metadata,
		Data:     make([]*Commentary, 0),
	}
}

func commentaryLoop(event *Event) {
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
			kafka.Publish(event.Metadata, commentary)
		}

	}
}
