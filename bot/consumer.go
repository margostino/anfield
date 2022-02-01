package bot

import (
	"fmt"
	"github.com/margostino/anfield/db"
	"log"
)

func (a App) consume() {
	user := <-a.subscriptions
	filter := db.GetUserFilter(user.Id)
	update := db.GetUpdateUser(user)
	document := a.db.Users.UpsertUser(filter, update)
	message := fmt.Sprintf("New Subscription from %s", document.Username)
	log.Println(message)
}
