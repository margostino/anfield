package processor

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"time"
)

// TODO: evaluate 2 consumers (goroutines): metadata + commentary
// TODO: calculate stats, bot sender
// This aggregation in consumer should happen once by URL/Event
func consume() {
	metadata := <-metadataBuffer
	event := &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
	commentaryLoop(event)
}

func commentaryLoop(event *domain.Event) {
	for {
		commentary := <-commentaryBuffer
		event.Data = append(event.Data, commentary)

		time.Sleep(100 * time.Millisecond)

		if commentary.Time == "" && commentary.Comment != "" {
			fmt.Printf("# %s\n", commentary.Comment)
		} else {
			fmt.Printf("%s - %s\n", commentary.Time, commentary.Comment)
		}

		//source.WriteOnFileIfUpdate(event)
	}
}
