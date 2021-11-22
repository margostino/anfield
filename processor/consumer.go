package processor

import (
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/io"
	"github.com/margostino/anfield/kafka"
	"time"
)

// TODO: calculate stats, bot sender
// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
// This aggregation in consumer should happen once by URL/Event
func consume(url string) {
	var event *domain.Event
	var metadata *domain.Metadata

	select {
	case metadata = <-metadataBuffer[url]:
		event = NewEvent(metadata)
	case <-time.After(10000 * time.Millisecond):
		metadata = &domain.Metadata{
			Url:      url,
			H2H:      "",
			Date:     "",
			HomeTeam: nil,
			AwayTeam: nil,
		}
	}

	event = NewEvent(metadata)
	enrich(event)
	save(event)
	done(url)
}

func NewEvent(metadata *domain.Metadata) *domain.Event {
	return &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
}

func enrich(event *domain.Event) {
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
			kafka.Produce(event.Metadata, commentary)
		}

	}
}

func save(event *domain.Event) {
	eventLines := toString(event)
	io.WriteOnFileIfUpdate(eventLines)
}

func done(url string) {
	waitGroups[url].Done()
}
