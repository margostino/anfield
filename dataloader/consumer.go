package dataloader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"time"
)

type Data struct {
	Comments []domain.Commentary
}

type Document struct {
	Metadata *domain.Metadata
	Data     *Data
}

func (a App) Consume() error {
	for {
		message, err := a.kafka.ReadMessage()

		if err != nil {
			break
		}

		a.upsertCommentary(message)
		a.upsertAssets(message)
	}
	return nil // TODO: tbd
}

func decode(result *mongo2.SingleResult) *Document {
	var document Document
	result.Decode(&document)
	return &document
}

func (a App) upsertAssets(message *domain.Message) {
	scores := a.scorer.CalculateScoring(message.Metadata.Lineups, message.Data)

	for key, value := range scores {
		filter := getAssetsFilter(key)
		update := getUpdateAssets(key, value)
		a.db.Assets.Upsert(filter, update)
		print(1)
	}

}

func (a App) upsertCommentary(message *domain.Message) {
	filter := getCommentaryFilter(message)
	update := getUpdateCommentary(message)
	result := a.db.Matches.Upsert(filter, update)
	document := decode(result)
	logging(document)
}

func getCommentaryFilter(message *domain.Message) bson.M {
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

func logging(document *Document) {
	id := common.GenerateEventID(document.Metadata.Url)
	dataLength := len(document.Data.Comments)
	message := fmt.Sprintf("New Message from %s with data length %d", id, dataLength)
	log.Println(message)
}
