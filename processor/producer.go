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


func metadata(url string) {
	metadata := scrapper.GetMetadata(url)
	metadataBuffer[url] <- metadata
}

// TODO: implement proper stop in loop but scan all partial events
func commentary(url string) {
	countDown := 0
	endOfEvent := false
	matchInProgress := true
	sent := 0
	eventName := strings.Split(url, "/")[7]
	stopFlag := context.Config().Realtime.StopFlag
	graceEndTime := context.Config().Realtime.GraceEndTime
	commentaryUrl := url + context.Config().Commentary.Params

	fmt.Printf("======== START: %s ========\n", eventName)

	for ok := true; ok; ok = matchInProgress {
		if endOfEvent && countDown == 0 {
			time.Sleep(graceEndTime * time.Millisecond)
			countDown += 1
		} else if endOfEvent && countDown == context.Config().Realtime.CountDown {
			matchInProgress = false
			break
		}
		rawEvents := scrapper.GetEvents(commentaryUrl)
		commentaries := normalize(*rawEvents)
		if sent != len(commentaries) {
			for _, commentary := range commentaries {
				newHash := common.Hash(commentary.Comment)
				if !common.InSlice(sent, newHash) {
					commentaryBuffer[url] <- commentary
					sent += 1
					if strings.Contains(commentary.Comment, stopFlag) {
						endOfEvent = true
					}
				}
			}
		}
	}

	fmt.Printf("======== END: %s ========\n", eventName)

	close(commentaryBuffer[url])
	asyncPool[url] <- true
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
