package processor

import (
	"fmt"
	"github.com/margostino/anfield/context"
	"github.com/margostino/anfield/scrapper"
	"github.com/margostino/anfield/source"
	"regexp"
)

func Process(urls []string) {
	config := context.Config()
	for _, url := range urls {
		infoUrl := url + config.Matches.InfoUrlParam
		date := scrapper.GetEventDate(infoUrl)
		commentsUrl := url + config.Matches.CommentUrlParam
		eventName, events := scrapper.GetEvents(commentsUrl)
		fmt.Printf("======== START: %s ========\n", eventName)
		save(date, events)
		fmt.Printf("======== END: %s ========\n", eventName)
	}

}

func save(eventDate string, events *[]string) {
	var time string
	for _, event := range *events {
		isTime, _ := regexp.MatchString("([0-9]?'|[0-9]{2}'|[0-9]{2}\\+[0-9]+')$", event)

		if isTime {
			time = event
		} else {
			line := fmt.Sprintf("%s;%s;%s\n", eventDate, time, event)
			fmt.Println(line)
			source.WriteOnFile(line)
			time = ""
		}
	}
}
