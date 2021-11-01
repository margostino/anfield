package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/scrapper"
	"github.com/margostino/anfield/source"
	"strings"
	"time"
)

var done = make(chan bool)
var msgs = make(chan string)

func Process(urls []string) {
	config := context.Config()
	for _, url := range urls {
		infoUrl := url + config.Matches.InfoUrlParam
		date := scrapper.GetEventDate(infoUrl)
		commentsUrl := url + config.Matches.CommentUrlParam
		go produce(date, commentsUrl)
		go consume()
		<-done
	}

}

// TODO: implement proper stop in loop but scan all partial events
func produce(date string, url string) {
	sent := make([]uint32, 0)
	matchInProgress := true
	endOfEvent := false
	stopFlag := context.Config().Realtime.StopFlag
	graceEndTime := context.Config().Realtime.GraceEndTime
	countDown := 0
	eventName := strings.Split(url, "/")[6]
	fmt.Printf("======== START: %s ========\n", eventName)
	for ok := true; ok; ok = matchInProgress {
		events := scrapper.GetEvents(date, url)
		common.Reverse(events)
		for _, event := range *events {
			newHash := common.Hash(event)
			if !common.Contains(sent, newHash) {
				msgs <- event
				sent = append(sent, newHash)
				if strings.Contains(event, stopFlag) {
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

func consume() {
	for {
		event := <-msgs
		time.Sleep(100 * time.Millisecond)
		fmt.Println(event)
		source.WriteOnFileIfUpdate(event)
	}
}
