package db

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func FilterBy(key string) bson.M {
	id := common.HashFrom(key)
	return bson.M{"_id": id}
}

func FilterPatternBy(key string, value string) bson.M {
	return bson.M{
		key: bson.M{
			"$regex": fmt.Sprintf(".*%s*", value), // TODO: improve regex, more robust
		},
	}
}

func GetAssetsPatternFilter(key string) bson.M {
	return bson.M{
		"name": bson.M{
			"$regex": fmt.Sprintf(".*%s*", key), // TODO: improve regex, more robust
		},
	}
}

func UpdateAssetQuery(asset *domain.Asset) bson.M {
	return bson.M{
		"$inc": bson.M{"score": asset.Score},
		"$set": bson.M{"name": asset.Name, "last_updated": asset.LastUpdated},
	}
}

func UpdateWalletQuery(value float64, assets []domain.WalletAssetDocument) bson.M {
	return bson.M{
		"$inc": bson.M{"wallet.budget": value},
		"$set": bson.M{"wallet.assets": assets, "wallet.last_updated": time.Now().UTC()},
	}
}

func InsertUserQuery(user *domain.User) bson.M {
	return bson.M{
		"_id":        common.HashFrom(string(user.SocialId)),
		"social_id":  user.SocialId,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"wallet":     user.Wallet,
	}
}

func GetInsertTransaction(transaction domain.Transaction) bson.M {
	return bson.M{
		"user_id":   transaction.UserId,
		"asset_id":  transaction.AssetId,
		"units":     transaction.Units,
		"value":     transaction.Value,
		"operation": transaction.Operation,
		"timestamp": transaction.Timestamp,
	}
}

func UpdateCompletionQuery(match *domain.Match) bson.M {
	return bson.M{
		"$set": bson.M{"metadata.finished": match.Metadata.Finished},
	}
}

func UpdateMatchQuery(match *domain.Match) bson.M {
	return bson.M{
		"$push": bson.M{"data.comments": match.Data},
		"$set":  bson.M{"metadata": match.Metadata, "lineups": match.Lineups},
	}
}
