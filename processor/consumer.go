package processor

import (
	"fmt"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: calculate stats, bot sender
// TODO: consumer does not need be a goroutine if it implements a infinite loop, unless we want extra process after that.
// This aggregation in consumer should happen once by URL/Event
func consume(url string) {
	var event *domain.Event
	var metadata *domain.Metadata
	timeout := configuration.ChannelTimeout()

	select {
	case metadata = <-metadataBuffer[url]:
		event = NewEvent(metadata)
	case <-time.After(timeout * time.Millisecond):
		log.Println("No metadata for", url)
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
	//h2h := event.Metadata.H2H
	for {
		commentary := <-commentaryBuffer[url]
		event.Data = append(event.Data, commentary)

		time.Sleep(100 * time.Millisecond) // TODO: configurable

		if end(commentary) {
			break
		} else {
			//printCommentary(h2h, commentary)
			loggingState(url, commentary)
			kafka.Produce(event.Metadata, commentary)
		}

	}
}

func loggingState(url string, commentary *domain.Commentary) {
	var time, additionalTime, totalTime int
	var completion float64

	if isTimedComment(commentary) {
		rawTime := strings.ReplaceAll(commentary.Time, "'", "")
		fullTime := strings.Split(rawTime, "+")
		time, _ = strconv.Atoi(fullTime[0])

		if len(fullTime) > 1 {
			additionalTime, _ = strconv.Atoi(fullTime[1])
			totalTime = time
		} else {
			totalTime = time + additionalTime
		}

		if totalTime > 90 {
			completion = 100
		} else {
			completion = float64(totalTime) * 100 / 90
		}

		message := fmt.Sprintf("[%s] > %.2f", url, completion)
		log.Println(message)
	}
}

func done(url string) {
	//waitGroups[url].Done()
	wg, _ := waitGroups.Load(url)
	wg.(*sync.WaitGroup).Done()
}
