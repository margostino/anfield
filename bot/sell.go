package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"regexp"
)

// TODO: define transaction output
// TODO: evaluate async operation and return a promise
// TODO: define retries

type SellInstruction struct {
	Command string
}

type Sell struct {
	Command      string
	Regex        *regexp.Regexp
	Db           *db.Database
	Users        *db.Collection
	Assets       *db.Collection
	Transactions *db.Collection
}

func (s SellInstruction) shouldReply(input string) bool {
	return s.Command == input
}

func (s SellInstruction) reply(_ *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	reply := "✏️   Please send Asset Name and Value separated by space.\n💡   Example:  salah 2"
	return nil, reply, true
}

func (s Sell) shouldReply(input string) bool {
	return s.Regex.MatchString(input)
}

func (s Sell) reply(update *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	var reply string
	input := update.Message.Text
	userId := update.Message.From.ID
	assetName, units := valuesFromBuySellOperation(input)
	asset, user, err := s.getTransactionParams(assetName, userId)
	ctx := context.Background()

	if err != nil {
		return nil, failureReply(), false
	}

	total := asset.Score * float64(units)
	newTransaction, transactionHandler := s.transaction(user, asset, units, total)

	session, err := s.Db.Client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, transactionHandler)

	if err != nil {
		log.Println("Error in transaction. Operation was aborted.", err)
		reply = failureReply()
	} else {
		reply = fmt.Sprintf("🤝   Successful transaction.\n"+
			"⏱️   Timestamp: %s\n"+
			"🛒   Operation: %s\n"+
			"⚽   AssetID: %s\n"+
			"🗳️   Units: %d\n"+
			"💵   Value: %.2f\n"+
			"💰   Total: %.2f",
			common.UTC(newTransaction.Timestamp),
			newTransaction.Operation,
			newTransaction.AssetId,
			newTransaction.Units,
			newTransaction.Value,
			-1*total,
		)
	}

	return nil, reply, false
}

func (s Sell) getTransactionParams(assetName string, userId int) (*domain.AssetDocument, *domain.UserDocument, error) {
	asset, err := s.getAsset(assetName)

	if err != nil {
		return nil, nil, err
	}

	user, err := s.getUser(userId)

	if err != nil {
		return nil, nil, err
	}

	return asset, user, nil
}

// TODO: tbd more than one asset results
func (s Sell) getAsset(name string) (*domain.AssetDocument, error) {
	var asset domain.AssetDocument
	filter := db.FilterPatternBy("name", name)
	err := s.Assets.FindOne(filter, &asset)
	if err != nil {
		log.Println(fmt.Sprintf("Asset %s search failed with error %s", name, err.Error()))
	}
	return &asset, err
}

func (s Sell) getUser(userId int) (*domain.UserDocument, error) {
	var user domain.UserDocument
	filter := db.FilterBy(string(userId))
	err := s.Users.FindOne(filter, &user)
	if err != nil {
		log.Println(fmt.Sprintf("User %d search failed with error %s", userId, err.Error()))
	}
	return &user, err
}

// TODO: validate user not found
func (s Sell) updateWallet(id string, value float64, context mongo.SessionContext) error {
	filter, update := db.UpdateWallet(id, value, nil)
	err := s.Users.Update(filter, update)
	if err != nil {
		log.Println(fmt.Sprintf("Wallet update failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("Wallet updated for user %s with budget %f", id, value))
	}
	return err
}

func (s Sell) insertTransaction(transaction *domain.Transaction, context mongo.SessionContext) error {
	document := db.GetInsertTransaction(transaction)
	err := s.Transactions.InsertWithContext(document, context)
	if err != nil {
		log.Println(fmt.Sprintf("Transaction failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("New Transaction from %s buying asset %s", transaction.UserId, transaction.AssetId))
	}
	return err
}

func (s Sell) transaction(user *domain.UserDocument, asset *domain.AssetDocument, units int, total float64) (*domain.Transaction, func(sessCtx mongo.SessionContext) (interface{}, error)) {
	newTransaction := &domain.Transaction{
		UserId:    user.Id,
		AssetId:   asset.Id,
		Value:     asset.Score,
		Units:     units,
		Operation: domain.BUY,
		Timestamp: common.Now(),
	}
	fc := func(context mongo.SessionContext) (interface{}, error) {
		if err := s.insertTransaction(newTransaction, context); err != nil {
			return nil, err
		}
		//return nil, errors.New("testing")
		if err := s.updateWallet(user.Id, total, context); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return newTransaction, fc
}