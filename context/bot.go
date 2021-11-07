package context

import (
	"github.com/margostino/anfield/configuration"
	"strings"
)

var subscriptions = make(map[string]string)
var matches = make([]string, 0)

func Subscriptions() map[string]string {
	return subscriptions
}

func Matches() []string {
	return matches
}

func Initialize() {
	for _, match := range configuration.Realtime().Matches {
		matches = append(matches, strings.Split(match, "/")[1])
	}
}

func Subscribe(username string, eventId string) {
	subscriptions[username] = eventId
}
