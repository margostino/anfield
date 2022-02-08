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
	"strconv"
	"strings"
)

// TODO: define transaction output
// TODO: evaluate async operation and return a promise
// TODO: define retries

type BuyInstruction struct {
	Command string
}

type Buy struct {
	Command      string
	Regex        *regexp.Regexp
	Db           *db.Database
	Users        *db.Collection
	Assets       *db.Collection
	Transactions *db.Collection
}

func (b BuyInstruction) shouldReply(input string) bool {
	return b.Command == input
}

func (b BuyInstruction) reply(_ *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	reply := "✏️   Please send Asset Name and Value separated by space.\n💡   Example:  salah 2"
	return nil, reply, true
}

func (b Buy) shouldReply(input string) bool {
	return b.Regex.MatchString(input)
}

func (b Buy) reply(update *tgbotapi.Update) (interface{}, string, bufferEnabled) {
	var reply string
	input := update.Message.Text
	userId := update.Message.From.ID
	assetName, units := extractValuesFrom(input)
	asset, user, err := b.getTransactionParams(assetName, userId)
	ctx := context.Background()

	if err != nil {
		return nil, failureReply(), false
	}

	total := -1 * asset.Score * float64(units)
	newTransaction, transactionHandler := b.transaction(user, asset, units, total)

	session, err := b.Db.Client.StartSession()
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

func failureReply() string {
	return fmt.Sprintf("🔴   Transaction can not be executed now.\n🙏🏻   Please try later.")
}

func (b Buy) getTransactionParams(assetName string, userId int) (*domain.AssetDocument, *domain.UserDocument, error) {
	asset, err := b.getAsset(assetName)

	if err != nil {
		return nil, nil, err
	}

	user, err := b.getUser(userId)

	if err != nil {
		return nil, nil, err
	}

	return asset, user, nil
}

func extractValuesFrom(input string) (string, int) {
	values := strings.Split(input, " ")
	assetName := values[0]
	units, _ := strconv.Atoi(values[1])
	return assetName, units
}

// TODO: tbd more than one asset results
func (b Buy) getAsset(name string) (*domain.AssetDocument, error) {
	var asset domain.AssetDocument
	filter := db.FilterPatternBy("name", name)
	err := b.Assets.FindOne(filter, &asset)
	if err != nil {
		log.Println(fmt.Sprintf("Asset %s search failed with error %s", name, err.Error()))
	}
	return &asset, err
}

func (b Buy) getUser(userId int) (*domain.UserDocument, error) {
	var user domain.UserDocument
	filter := db.FilterBy(string(userId))
	err := b.Users.FindOne(filter, &user)
	if err != nil {
		log.Println(fmt.Sprintf("User %d search failed with error %s", userId, err.Error()))
	}
	return &user, err
}

// TODO: validate user not found
func (b Buy) updateWallet(id string, value float64, context mongo.SessionContext) error {
	var wallet domain.WalletDocument
	filter, update := db.UpdateWallet(id, value)
	err := b.Users.Update(filter, update, &wallet)
	if err != nil {
		log.Println(fmt.Sprintf("Wallet update failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("Wallet updated for user %s with budget %f", id, value))
	}
	return err
}

func (b Buy) insertTransaction(transaction *domain.Transaction, context mongo.SessionContext) error {
	document := db.GetInsertTransaction(transaction)
	err := b.Transactions.InsertWithContext(document, context)
	if err != nil {
		log.Println(fmt.Sprintf("Transaction failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("New Transaction from %s buying asset %s", transaction.UserId, transaction.AssetId))
	}
	return err
}

func (b Buy) transaction(user *domain.UserDocument, asset *domain.AssetDocument, units int, total float64) (*domain.Transaction, func(sessCtx mongo.SessionContext) (interface{}, error)) {
	newTransaction := &domain.Transaction{
		UserId:    user.Id,
		AssetId:   asset.Id,
		Value:     asset.Score,
		Units:     units,
		Operation: domain.BUY,
		Timestamp: common.Now(),
	}
	fc := func(context mongo.SessionContext) (interface{}, error) {
		if err := b.insertTransaction(newTransaction, context); err != nil {
			return nil, err
		}
		//return nil, errors.New("testing")
		if err := b.updateWallet(user.Id, total, context); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return newTransaction, fc
}
