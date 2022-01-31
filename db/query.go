package db

import (
	"encoding/hex"
	"github.com/margostino/anfield/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUrlFilter(url string) bson.M {
	_, _, identifier := common.ExtractTeamsFrom(url)
	hex := hex.EncodeToString([]byte(identifier + identifier))
	id, err := primitive.ObjectIDFromHex(hex)
	common.Check(err)
	return bson.M{"_id": id}
}
