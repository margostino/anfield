package bot

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
)

type Trade struct {
	Total         float64
	Transaction   domain.Transaction
	WalletAssets  []domain.WalletAssetDocument
	AtomicHandler func(context mongo.SessionContext) (interface{}, error)
}

func valuesTradeOperation(input string) (string, int) {
	values := strings.Split(input, " ")
	assetName := values[0]
	units, _ := strconv.Atoi(values[1])
	return assetName, units
}

func createTrade(operation string, user *domain.UserDocument, asset *domain.AssetDocument, units int) Trade {
	var total float64

	if operation == domain.SELL {
		units *= -1
		total = asset.Score * float64(units)
	} else {
		total = -1 * asset.Score * float64(units)
	}

	transaction := domain.Transaction{
		UserId:    user.Id,
		AssetId:   asset.Id,
		Value:     asset.Score,
		Units:     units,
		Operation: operation,
		Timestamp: common.Now(),
	}

	walletAssets := findAndUpdateAssetInWallet(operation, asset, units, user.Wallet.Assets)

	return Trade{
		Total:        total,
		Transaction:  transaction,
		WalletAssets: walletAssets,
	}
}

// TODO: implement an upsert operation pushing new Asset in Wallet to avoid rewrite every time the array
func findAndUpdateAssetInWallet(operation string, currentAsset *domain.AssetDocument, units int, assets []domain.WalletAssetDocument) []domain.WalletAssetDocument {
	var multiplier int
	var exists = false
	updatedAssets := make([]domain.WalletAssetDocument, 0)

	if operation == domain.SELL {
		multiplier = -1
	} else {
		multiplier = 1
	}

	for _, asset := range assets {
		if asset.Id == currentAsset.Id {
			asset.Units += multiplier * units
			asset.Value = currentAsset.Score
			exists = true
		}

		if asset.Units > 0 {
			updatedAssets = append(updatedAssets, asset)
		} else {
			exists = false
		}

	}

	if !exists && operation == domain.BUY {
		asset := domain.WalletAssetDocument{
			Id:    currentAsset.Id,
			Units: units,
			Value: currentAsset.Score,
		}
		updatedAssets = append(updatedAssets, asset)
	}

	return updatedAssets
}
