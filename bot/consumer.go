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
	wallet := &domain.Wallet{
		Budget:      100000,
		LastUpdated: time.Now(),
	}
	a.createUser(user, wallet)
}

func (a App) createUser(user domain.User, wallet *domain.Wallet) {
	update := db.GetInsertUser(user, wallet)
	err := a.db.Users.InsertUser(update)

	if err == nil {
		message := fmt.Sprintf("New Subscription from %s", user.Username)
		log.Println(message)
	}

}
