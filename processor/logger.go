package processor

import (
	"fmt"
	"github.com/margostino/anfield/configuration"
	"github.com/margostino/anfield/domain"
	"log"
	"strconv"
	"strings"
	"sync"
)

var state = sync.Map{}

type Logger struct {
	configuration *configuration.Configuration
}

func NewLogger(configuration *configuration.Configuration) *Logger {
	return &Logger{
		configuration: configuration,
	}
}

func (l *Logger) log(url string, commentary *domain.Commentary) {
	step := l.configuration.Logger.CompletionStep
	var time, additionalTime, totalTime int
	var completionFloat float64

	event := strings.Split(url, "/")[7]

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
			completionFloat = 100
		} else {
			completionFloat = float64(totalTime) * 100 / 90
		}

		completion := int(completionFloat)

		value, _ := state.Load(event)
		if (completion == 1 || completion%step == 0) && (value == nil || value.(int) < completion) {
			state.Store(event, completion)
			message := fmt.Sprintf("[%s] ==> %d%%", event, completion)
			log.Println(message)
		}

	}
}
