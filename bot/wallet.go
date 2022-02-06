package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
)

type Wallet struct {
	Command string
	Users   *db.Collection
}

func (w Wallet) shouldReply(input string) bool {
	return w.Command == input
}

func (w Wallet) reply(update *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	userId := update.Message.From.ID
	wallet, err := w.getWallet(userId)
	if err != nil {
		return nil, "Wallet not found", false
	}

	reply := fmt.Sprintf("💰   Budget: $ %.2f.\n"+
		"📅   Last updated: %s",
		wallet.Budget,
		common.UTC(wallet.LastUpdated))

	return nil, reply, false
}

func (w Wallet) getWallet(userId int) (*domain.WalletDocument, error) {
	var user domain.UserDocument
	filter := db.FilterBy(string(userId))
	err := w.Users.FindOne(filter, &user)
	if err != nil {
		log.Println(fmt.Sprintf("User %d search failed with error %s", userId, err.Error()))
	}
	return user.Wallet, err
}
