package db

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

func GetUrlFilter(url string) bson.M {
	_, _, identifier := common.ExtractTeamsFrom(url)
	hex := hex.EncodeToString([]byte(identifier + identifier))
	id, err := primitive.ObjectIDFromHex(hex)
	common.Check(err)
	return bson.M{"_id": id}
}

func GetUserFilter(userId int64) bson.M {
	id := hashFrom(string(userId))
	return bson.M{"_id": id}
}

func GetAssetsFilter(key string) bson.M {
	id := hashFrom(key)
	return bson.M{"_id": id}
}

func GetAssetsPatternFilter(key string) bson.M {
	return bson.M{
		"name": bson.M{
			"$regex": fmt.Sprintf(".*%s*", key), // TODO: improve regex, more robust
		},
	}
}

func GetUpdateAssets(name string, score float64) bson.M {
	return bson.M{
		"$inc": bson.M{"score": score},
		"$set": bson.M{"name": name, "last_updated": time.Now().UTC()},
	}
}

func GetUpdateUser(budget float64) bson.M {
	return bson.M{
		"$inc": bson.M{"wallet.budget": budget},
		"$set": bson.M{"wallet.last_updated": time.Now().UTC()},
	}
}

func GetInsertUser(user domain.User, wallet *domain.Wallet) bson.M {
	return bson.M{
		"_id":        hashFrom(string(user.Id)),
		"social_id":  user.Id,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"wallet":     wallet,
	}
}

func GetInsertTransaction(transaction *domain.Transaction) bson.M {
	return bson.M{
		"user_id":   transaction.UserId,
		"asset_id":  transaction.AssetId,
		"units":     transaction.Units,
		"value":     transaction.Value,
		"timestamp": transaction.Timestamp,
	}
}

func GetUpdateCompletion(message *domain.Message) bson.M {
	return bson.M{
		"$set": bson.M{"metadata.finished": message.Metadata.Finished},
	}
}

func GetUpdateCommentary(message *domain.Message) bson.M {
	return bson.M{
		"$push": bson.M{"data.comments": message.Data},
		"$set":  bson.M{"metadata": message.Metadata, "lineups": message.Lineups},
	}
}

func hashFrom(key string) string {
	hash := sha1.New()
	io.WriteString(hash, key)
	return hex.EncodeToString(hash.Sum(nil))
}
