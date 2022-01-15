package dataloader

import (
	"fmt"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"log"
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
		a.upsertScoring(message)
	}
	return nil // TODO: tbd
}

func decode(result *mongo2.SingleResult) *Document {
	var document Document
	result.Decode(&document)
	return &document
}

func (a App) upsertScoring(message *domain.Message) {
	scores := a.scorer.CalculateScoring(message.Metadata.Lineups, message.Data)

	for key, value := range scores {
		filter := bson.M{"player": key}
		update := bson.M{
			"$push": bson.M{"scores": value},
			"$set":  bson.M{"player": key},
		}
		a.db.Matches.Upsert(filter, update)
	}

}

func (a App) upsertCommentary(message *domain.Message) {
	filter := getFilter(message)
	update := getUpdateDoc(message)
	result := a.db.Matches.Upsert(filter, update)
	document := decode(result)
	logging(document)
}

func getFilter(message *domain.Message) bson.M {
	return bson.M{"metadata.url": message.Metadata.Url}
}

func getUpdateDoc(message *domain.Message) bson.M {
	return bson.M{
		"$push": bson.M{"data.comments": message.Data},
		"$set":  bson.M{"metadata": message.Metadata},
	}
}

func logging(document *Document) {
	id := document.Metadata.Id
	dataLength := len(document.Data.Comments)
	message := fmt.Sprintf("New Message from %s with data length %d", id, dataLength)
	log.Println(message)
}
