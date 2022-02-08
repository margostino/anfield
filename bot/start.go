package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
)

// TODO: enrich stats with more information (last update, highest/lowest, etc...)
// TODO: add command to explain stats
// TODO: add command to alert trends and changes and automate buy/sell operation given conditions (e.g. threshold)
// TODO: persist stats history and matches reply (realtime + batch contribution to avoid duplications)

type Start struct {
	Command string
	Users   *db.Collection
}

func (s Start) shouldReply(input string) bool {
	return s.Command == input
}

func (s Start) reply(update *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	var name string
	user := User(update)
	if user.FirstName != "" {
		name = user.FirstName
	} else {
		name = user.Username
	}

	go s.subscribe(user)

	return nil, fmt.Sprintf("👋   Hi %s, Welcome to Anfield!", name), false
}

func (s Start) subscribe(user *domain.User) {
	var message string
	document := db.InsertUserQuery(user)
	err := s.Users.Insert(document)

	if err != nil && db.IsDuplicatedWrite(err) {
		message = fmt.Sprintf("User %s already exists", user.Username)
	} else {
		message = fmt.Sprintf("New Subscription from %s", user.Username)
	}
	log.Println(message)
}

func User(update *tgbotapi.Update) *domain.User {
	from := update.Message.From
	return &domain.User{
		SocialId:  from.ID,
		Username:  from.UserName,
		FirstName: from.FirstName,
		LastName:  from.LastName,
		Wallet: &domain.Wallet{
			Budget:      100000,
			LastUpdated: common.Now(),
		},
	}
}
