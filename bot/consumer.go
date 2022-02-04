package bot

import (
	"fmt"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
	"time"
)

func (a App) consume() {
	user := <-a.subscriptions
	user.Wallet = &domain.Wallet{
		Budget:      100000,
		LastUpdated: time.Now().UTC(),
	}
	a.createUser(&user)
}

func (a App) createUser(user *domain.User) {
	document := db.InsertUserQuery(user)
	err := a.db.Users.Insert(document)

	if err == nil {
		message := fmt.Sprintf("New Subscription from %s", user.Username)
		log.Println(message)
	}

}
