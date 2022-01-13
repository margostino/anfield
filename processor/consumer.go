package processor

import (
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	"log"
	"time"
)

// TODO: calculate stats, bot sender
// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
// This aggregation in consumer should happen once by URL/Event
func (a App) consume(url string) {
	var event *domain.Event
	var metadata *domain.Metadata
	var consumedLineups *domain.Lineups
	var consumedDate string

	timeout := a.configuration.ChannelTimeout()

	select {
	case consumedDate = <-a.channels.matchDate[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No metadata for", url)
		consumedDate = ""
	}

	select {
	case consumedLineups = <-a.channels.lineups[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No lineups for", url)
		consumedLineups = nil
	}

	metadata = &domain.Metadata{
		Url:     url,
		H2H:     "", // TODO: generate ID //h2h := fmt.Sprintf("%s vs %s", homeTeam.Name, awayTeam.Name)
		Lineups: consumedLineups,
		Date:    consumedDate,
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
			kafka.Produce(event.Metadata, nil)
		} else {
			//logging(url, commentary)
			kafka.Produce(event.Metadata, commentary)
		}
	}
}

func done(url string) {
	//wg, _ := waitGroups.Load(url)
	//wg.(*sync.WaitGroup).Done()
	waitGroup.Done()
}
