package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/db"
	"log"
	"regexp"
)

type bufferEnabled bool

type Action interface {
	shouldReply(input string) bool
	reply(update *tgbotapi.Update) (interface{}, string, bufferEnabled)
}

// NewActions TODO: tbd by config
func NewActions(db *db.Database) []Action {
	var actions = make([]Action, 0)
	start := StartAction(db.Users)
	wallet := WalletAction(db.Users)
	buyInstruction, buy := BuyAction(db)
	actions = append(actions, start, wallet, buyInstruction, buy)
	return actions
}

func StartAction(users *db.Collection) Start {
	return Start{
		Command: "/start",
		Users:   users,
	}
}

func WalletAction(users *db.Collection) Wallet {
	return Wallet{
		Command: "/wallet",
		Users:   users,
	}
}

func BuyAction(db *db.Database) (BuyInstruction, Buy) {
	instruction := BuyInstruction{
		Command: "/buy",
	}
	buyCommand := "^\\/buy [A-Za-z]+ [1-9]+[0-9]*$"
	regex, err := regexp.Compile(buyCommand)

	if err != nil {
		log.Println("Error compiling Regex for Buy Action", err)
	}

	buy := Buy{
		Command:      buyCommand,
		Regex:        regex,
		Db:           db,
		Users:        db.Users,
		Assets:       db.Assets,
		Transactions: db.Transactions,
	}

	return instruction, buy
}
