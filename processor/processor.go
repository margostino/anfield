package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scrapper"
	"strings"
	"sync"
	"time"
)

var done = make(chan bool)
var msgs = make(chan domain.Commentary)

// Process TODO: spawn one process per URL
func Process(urls []string) {
	wg := common.WaitGroup(len(urls))
	for _, url := range urls {
		go scan(url, wg)
	}
	wg.Wait()
}

func scan(url string, waitGroup *sync.WaitGroup) {
	config := context.Config()
	metadata := scrapper.GetMetadata(url)
	commentaryUrl := url + config.Commentary.Params
	event := &domain.Event{
		Metadata: metadata,
		Data:     make([]*domain.Commentary, 0),
	}
	go produce(commentaryUrl)
	go consume(event)
	<-done
	waitGroup.Done()
}

// TODO: implement proper stop in loop but scan all partial events
func produce(url string) {
	countDown := 0
	endOfEvent := false
	matchInProgress := true
	sent := 0
	eventName := strings.Split(url, "/")[7]
	stopFlag := context.Config().Realtime.StopFlag
	graceEndTime := context.Config().Realtime.GraceEndTime

	fmt.Printf("======== START: %s ========\n", eventName)

	for ok := true; ok; ok = matchInProgress {
		if endOfEvent && countDown == 0 {
			time.Sleep(graceEndTime * time.Millisecond)
			countDown += 1
		} else if endOfEvent && countDown == context.Config().Realtime.CountDown {
			matchInProgress = false
			break
		}
		rawEvents := scrapper.GetEvents(url)
		commentaries := normalize(*rawEvents)
		if sent != len(commentaries) {
			for _, commentary := range commentaries {
				newHash := common.Hash(commentary.Comment)
				if !common.InSlice(sent, newHash) {
					msgs <- *commentary
					sent += 1
					if strings.Contains(commentary.Comment, stopFlag) {
						endOfEvent = true
					}
				}
			}
		}
	}

	fmt.Printf("======== END: %s ========\n", eventName)

	close(msgs)
	done <- true
}

func consume(event *domain.Event) {
	for {
		msg := <-msgs
		time.Sleep(100 * time.Millisecond)
		event.Data = append(event.Data, &msg)

		if msg.Time == "" && msg.Comment != "" {
			fmt.Printf("# %s\n", msg.Comment)
		} else {
			fmt.Printf("%s - %s\n", msg.Time, msg.Comment)
		}

		//source.WriteOnFileIfUpdate(event)
	}
}

func normalize(comments []string) []*domain.Commentary {
	var time string
	var commentaries = make([]*domain.Commentary, 0)

	for _, value := range comments {
		if common.IsTimeCounter(value) {
			time = value
		} else {
			commentary := domain.Commentary{
				Time:    time,
				Comment: value,
			}
			commentaries = append(commentaries, &commentary)
			time = ""
		}
	}
	common.Reverse(&commentaries)
	return commentaries
}
