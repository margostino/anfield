package bot

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
)

// TODO: define transaction output
// TODO: evaluate async operation and return a promise
// TODO: define retries
func (a App) buy(userId int, assetName string, units int) (*domain.Transaction, error) {
	asset, err := a.getAsset(assetName)

	if err != nil {
		return nil, err
	}

	user, err := a.getUser(userId)

	if err != nil {
		return nil, err
	}

	total := -1 * asset.Score * float64(units)
	newTransaction := transaction(user, asset, units)

	err = a.insertTransaction(newTransaction)

	if err != nil {
		return newTransaction, err
	}

	err = a.updateWallet(user.Id, total)

	return newTransaction, err
}

func transaction(user *domain.UserDocument, asset *domain.AssetDocument, units int) *domain.Transaction {
	return &domain.Transaction{
		UserId:    user.Id,
		AssetId:   asset.Id,
		Value:     asset.Score,
		Units:     units,
		Operation: domain.BUY,
		Timestamp: common.Now(),
	}
}

// TODO: tbd more than one asset results
func (a App) getAsset(name string) (*domain.AssetDocument, error) {
	var asset domain.AssetDocument
	filter := db.FilterPatternBy("name", name)
	err := a.db.Assets.FindOne(filter, &asset)
	if err != nil {
		log.Println(fmt.Sprintf("Asset %s search failed with error %s", name, err.Error()))
	}
	return &asset, err
}

func (a App) getUser(userId int) (*domain.UserDocument, error) {
	var user domain.UserDocument
	filter := db.FilterBy(string(userId))
	err := a.db.Users.FindOne(filter, &user)
	if err != nil {
		log.Println(fmt.Sprintf("User %d search failed with error %s", userId, err.Error()))
	}
	return &user, err
}

// TODO: validate user not found
func (a App) updateWallet(id string, value float64) error {
	var document domain.WalletDocument
	filter, update := db.UpsertWallet(id, value)
	err := a.db.Users.Upsert(filter, update, &document)
	if err != nil {
		log.Println(fmt.Sprintf("Wallet update failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("Wallet updated for user %s with budget %f", id, value))
	}
	return err
}

func (a App) insertTransaction(transaction *domain.Transaction) error {
	document := db.GetInsertTransaction(transaction)
	err := a.db.Transactions.Insert(document)
	if err != nil {
		log.Println(fmt.Sprintf("Transaction failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("New Transaction from %s buying asset %s", transaction.UserId, transaction.AssetId))
	}
	return err
}
