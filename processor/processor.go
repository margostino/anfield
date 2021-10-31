package processor

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/scrapper"
	"github.com/margostino/anfield/source"
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
		//eventName := strings.Split(commentsUrl, "/")[7]
		//fmt.Printf("======== START: %s ========\n", eventName)
		go produce(date, commentsUrl)
		go consume()
		<-done
		//fmt.Printf("======== END: %s ========\n", eventName)
	}

}

// TODO: implement stop in loop but scan all partial events
func produce(date string, url string) {
	sent := make([]uint32, 0)
	matchInProgress := true

	for ok := true; ok; ok = matchInProgress {
		events := scrapper.GetEvents(date, url)
		common.Reverse(events)
		for _, event := range *events {
			newHash := common.Hash(event)
			if !common.Contains(sent, newHash) {
				// TODO: add stop the game
				msgs <- event
				sent = append(sent, newHash)
			}
		}
	}

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
