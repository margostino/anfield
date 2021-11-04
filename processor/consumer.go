package processor

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/source"
	"time"
)

// TODO: calculate stats, bot sender
// This aggregation in consumer should happen once by URL/Event
func consume(url string) {
	metadata := <-metadataBuffer[url]
	event := &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
	commentaryLoop(event)
	source.WriteOnFileIfUpdate(event)
	waitGroups[url].Done()
}

func commentaryLoop(event *domain.Event) {
	url := event.Metadata.Url
	h2h := event.Metadata.H2H
	for {
		commentary := <-commentaryBuffer[url]
		event.Data = append(event.Data, commentary)

		time.Sleep(100 * time.Millisecond)

		if end(commentary) {
			break
		} else if isTimedComment(commentary) {
			fmt.Printf("[%s] # %s\n", h2h, commentary.Comment)
		} else {
			fmt.Printf("[%s] # %s - %s\n", h2h, commentary.Time, commentary.Comment)
		}

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
