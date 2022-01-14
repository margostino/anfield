package processor

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"log"
	"time"
)

// TODO: calculate stats, bot sender
// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
// This aggregation in consumer should happen once by URL/Event
func (a App) consume(url string) {
	var event *domain.Event
	var metadata *domain.Metadata
	var lineups *domain.Lineups
	var date string

	timeout := a.configuration.ChannelTimeout()

	select {
	case date = <-a.channels.matchDate[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No metadata for", url)
		date = ""
	}

	select {
	case lineups = <-a.channels.lineups[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No lineups for", url)
		lineups = nil
	}

	metadata = &domain.Metadata{
		Url:     url,
		Id:      common.GenerateEventID(url),
		Lineups: lineups,
		Date:    date,
	}

	event = NewEvent(metadata)
	a.enrich(event)
	done(url)
}

func NewEvent(metadata *domain.Metadata) *domain.Event {
	return &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
}

func (a App) enrich(event *domain.Event) {
	url := event.Metadata.Url
	//h2h := event.Metadata.H2H
	for {
		commentary := <-a.channels.commentary[url]
		event.Data = append(event.Data, commentary)

		time.Sleep(100 * time.Millisecond) // TODO: configurable

		if end(commentary) {
			break
		} else if notStarted(commentary) {
			a.kafka.Produce(event.Metadata, nil)
		} else {
			a.logger.log(url, commentary)
			a.kafka.Produce(event.Metadata, commentary)
		}
	}
}

func done(url string) {
	//wg, _ := waitGroups.Load(url)
	//wg.(*sync.WaitGroup).Done()
	waitGroup.Done()
}
