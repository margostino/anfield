package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/scrapper"
	"strings"
	"time"
)

var done = make(chan bool)
var msgs = make(chan domain.Commentary)

// Process TODO: spawn one process per URL
func Process(urls []string) {
	config := context.Config()
	for _, url := range urls {
		metadata := scrapper.GetMetadata(url)
		commentaryUrl := url + config.Commentary.Params
		event := &domain.Event{
			Metadata: metadata,
			Data:     make([]*domain.Commentary, 0),
		}
		go produce(commentaryUrl)
		go consume(event)
		<-done
	}

}

// TODO: implement proper stop in loop but scan all partial events
func produce(url string) {
	sent := make([]uint32, 0)
	matchInProgress := true
	endOfEvent := false
	stopFlag := context.Config().Realtime.StopFlag
	graceEndTime := context.Config().Realtime.GraceEndTime
	countDown := 0
	eventName := strings.Split(url, "/")[7]
	fmt.Printf("======== START: %s ========\n", eventName)
	for ok := true; ok; ok = matchInProgress {
		rawEvents := scrapper.GetEvents(url)
		commentaries := normalize(*rawEvents)
		for _, commentary := range commentaries {
			newHash := common.Hash(commentary.Comment)
			if !common.InSlice(sent, newHash) {
				msgs <- *commentary
				sent = append(sent, newHash)
				if strings.Contains(commentary.Comment, stopFlag) {
					endOfEvent = true
				}
			} else {
				if endOfEvent && countDown == 0 {
					time.Sleep(graceEndTime * time.Millisecond)
					countDown += 1
				} else if endOfEvent && countDown == context.Config().Realtime.CountDown {
					matchInProgress = false
					break
				}
			}
		}
	}
	fmt.Printf("======== END: %s ========\n", eventName)
	fmt.Println("Before closing channel")
	close(msgs)
	fmt.Println("Before passing true to done")
	done <- true
}

func consume(event *domain.Event) {
	for {
		msg := <-msgs
		time.Sleep(100 * time.Millisecond)
		event.Data = append(event.Data, &msg)
		fmt.Println(msg.Comment)
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
