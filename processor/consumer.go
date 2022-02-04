package processor

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"log"
	"time"
)

// TODO: calculate stats, bot sender
// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
// This aggregation in consumer should happen once by URL/Event
func (a App) consume(url string) {
	var event *domain.Match
	var metadata *domain.Metadata
	var lineups *domain.Lineups
	var date time.Time

	timeout := a.configuration.ChannelTimeout()

	select {
	case date = <-a.channels.matchDate[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No metadata for", url)
		date = time.Time{}
	}

	select {
	case lineups = <-a.channels.lineups[url]:
	case <-time.After(timeout * time.Millisecond):
		log.Println("No lineups for", url)
		lineups = nil
	}

	metadata = &domain.Metadata{
		Url:  url,
		Date: date,
	}

	event = NewMatch(metadata, lineups)
	a.enrich(event)
	done()
}

func NewMatch(metadata *domain.Metadata, lineups *domain.Lineups) *domain.Match {
	return &domain.Match{
		Metadata: metadata,
		Lineups:  lineups,
		Data:     nil,
	}
}

func (a App) enrich(match *domain.Match) {
	url := match.Metadata.Url
	for {
		commentary := <-a.channels.commentary[url]
		match.Data = commentary

		//time.Sleep(100 * time.Millisecond) // TODO: configurable

		if end(commentary) { // TODO: define and set TTL just in case
			a.logger.info(fmt.Sprintf("End of match %s", url))
			match.Metadata.Finished = true
			a.kafka.Produce(match.Metadata, nil, nil)
			break
		} else if notStarted(commentary) {
			a.logger.info(fmt.Sprintf("Match %s is not started", url))
			a.kafka.Produce(match.Metadata, nil, nil)
		} else {
			a.logger.log(url, commentary)
			a.kafka.Produce(match.Metadata, commentary, match.Lineups)
		}
	}
}

func done() {
	//wg, _ := waitGroups.Load(url)
	//wg.(*sync.WaitGroup).Done()
	waitGroup.Done()
}
