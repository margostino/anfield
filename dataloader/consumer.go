package dataloader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"time"
)

func (a App) Consume() error {
	for {
		message, err := a.kafka.ReadMessage()

		if err != nil {
			break
		}

		if message.Metadata.Finished {
			a.upsertCompletion(message)
		} else {
			a.upsertCommentary(message)
			a.upsertAssets(message)
		}
	}
	return nil // TODO: tbd
}

func (a App) upsertAssets(message *domain.Message) {
	scores := a.scorer.CalculateScoring(message.Lineups, message.Data)

	for key, value := range scores {
		filter := getAssetsFilter(key)
		update := getUpdateAssets(key, value)
		a.db.Assets.Upsert(filter, update)
	}

}

func (a App) upsertCompletion(message *domain.Message) {
	filter := getUrlFilter(message)
	update := getUpdateCompletion(message)
	document := a.db.Matches.Upsert(filter, update)
	logging(document)
}

func (a App) upsertCommentary(message *domain.Message) {
	filter := getUrlFilter(message)
	update := getUpdateCommentary(message)
	document := a.db.Matches.Upsert(filter, update)
	logging(document)
}

func getUrlFilter(message *domain.Message) bson.M {
	//return bson.M{"metadata.url": message.Metadata.Url}
	_, _, identifier := common.ExtractTeamsFrom(message.Metadata.Url)
	hex := hex.EncodeToString([]byte(identifier + identifier))
	id, err := primitive.ObjectIDFromHex(hex)
	common.Check(err)
	return bson.M{"_id": id}
}

func getAssetsFilter(key string) bson.M {
	hash := sha1.New()
	io.WriteString(hash, key)
	id := hex.EncodeToString(hash.Sum(nil))
	return bson.M{"_id": id}
}

func getUpdateCompletion(message *domain.Message) bson.M {
	return bson.M{
		"$set": bson.M{"finished": message.Metadata.Finished},
	}
}

func getUpdateCommentary(message *domain.Message) bson.M {
	return bson.M{
		"$push": bson.M{"data.comments": message.Data},
		"$set":  bson.M{"metadata": message.Metadata},
	}
}

func getUpdateAssets(name string, score float64) bson.M {
	return bson.M{
		"$inc": bson.M{"score": score},
		"$set": bson.M{"name": name, "last_updated": time.Now().UTC()},
	}
}

func logging(document *domain.Document) {
	id := common.GenerateEventID(document.Metadata.Url)
	dataLength := len(document.Data.Comments)
	message := fmt.Sprintf("New Message from %s with data length %d", id, dataLength)
	log.Println(message)
}
