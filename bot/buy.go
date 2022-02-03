package bot

import (
	"fmt"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
	"time"
)

// TODO: define transaction output
// TODO: evaluate async operation and return a promise
func (a App) buy(userId int64, key string, units int) {
	asset := a.getAsset(key)
	user := a.getUser(userId)
	total := -1 * asset.Score * float64(units)
	transaction := newTransaction(user, asset, units)
	a.updateWallet(userId, total)
	a.insertTransaction(transaction)
}

func newTransaction(user *domain.UserDocument, asset *domain.AssetDocument, units int) *domain.Transaction {
	return &domain.Transaction{
		UserId:    user.Id,
		AssetId:   asset.Id,
		Value:     asset.Score,
		Units:     units,
		Timestamp: time.Now(),
	}
}

// TODO: validate asset not found
func (a App) getAsset(key string) *domain.AssetDocument {
	filter := db.GetAssetsPatternFilter(key)
	return a.db.Assets.FindOneAsset(filter)
}

func (a App) getUser(userId int64) *domain.UserDocument {
	filter := db.GetUserFilter(userId)
	return a.db.Users.FindOneUser(filter)
}

// TODO: validate user not found
func (a App) updateWallet(userId int64, total float64) {
	filter := db.GetUserFilter(userId)
	update := db.GetUpdateUser(total)
	a.db.Users.UpsertUser(filter, update)
}

func (a App) insertTransaction(transaction *domain.Transaction) {
	document := db.GetInsertTransaction(transaction)
	err := a.db.Transactions.Insert(document)

	if err == nil {
		message := fmt.Sprintf("New Transaction from %s buting asset %s", transaction.UserId, transaction.AssetId)
		log.Println(message)
	}
}
