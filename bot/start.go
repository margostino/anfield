package bot

import (
	"fmt"
	"github.com/margostino/anfield/domain"
)

// TODO: enrich stats with more information (last update, highest/lowest, etc...)
// TODO: add command to explain stats
// TODO: add command to alert trends and changes and automate buy/sell operation given conditions (e.g. threshold)
// TODO: persist stats history and matches reply (realtime + batch contribution to avoid duplications)

func shouldStart(message string) bool {
	return message == "/start"
}

func startReply(user *domain.User) (interface{}, string) {
	var name string

	if user.FirstName != "" {
		name = user.FirstName
	} else {
		name = user.Username
	}

	return nil, fmt.Sprintf("Hi %s, Welcome to Anfield!", name)
}

func (a App) subscribe(user *domain.User) {
	a.subscriptions <- *user
}