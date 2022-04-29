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
	assetName, units := valuesTradeOperation(input)
	asset, user, err := b.getTransactionParams(assetName, userId)
	ctx := context.Background()

	if err != nil {
		return nil, failureReply(), false
	}

	trade, atomicHandler := b.transaction(domain.BUY, user, asset, units)

	session, err := b.Db.Client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, atomicHandler)

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
			common.UTC(trade.Transaction.Timestamp),
			trade.Transaction.Operation,
			common.Mask(trade.Transaction.AssetId),
			trade.Transaction.Units,
			trade.Transaction.Value,
			-1*trade.Total,
		)
	}

	return nil, reply, false
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

func (b Buy) updateWallet(id string, budget float64, assets []domain.WalletAssetDocument, context mongo.SessionContext) error {
	filter, update := db.UpdateWallet(id, budget, assets)
	err := b.Users.UpdateWithContext(filter, update, context)
	if err != nil {
		log.Println(fmt.Sprintf("Wallet update failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("Wallet updated for user %s with budget %f", id, budget))
	}
	return err
}

func (b Buy) insertTransaction(transaction domain.Transaction, context mongo.SessionContext) error {
	document := db.GetInsertTransaction(transaction)
	err := b.Transactions.InsertWithContext(document, context)
	if err != nil {
		log.Println(fmt.Sprintf("Transaction failed with error %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("New Transaction from %s buying asset %s", transaction.UserId, transaction.AssetId))
	}
	return err
}

func (b Buy) transaction(operation string, user *domain.UserDocument, asset *domain.AssetDocument, units int) (Trade, func(sessCtx mongo.SessionContext) (interface{}, error)) {
	trade := createTrade(operation, user, asset, units)

	atomicHandler := func(context mongo.SessionContext) (interface{}, error) {
		if err := b.insertTransaction(trade.Transaction, context); err != nil {
			return nil, err
		}
		if err := b.updateWallet(trade.Transaction.UserId, trade.Total, trade.WalletAssets, context); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return trade, atomicHandler
}
