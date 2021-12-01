package dataloader

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Data struct {
	Comments []domain.Commentary
}

type Document struct {
	Metadata *domain.Metadata
	Data     *Data
}

func Consume() {
	for {
		message, err := kafka.ReadMessage()

		if err != nil {
			break
		}

		filter := bson.M{"metadata.url": message.Metadata.Url}
		opt := upsertOptions()
		update := bson.M{
			"$push": bson.M{"data.comments": message.Data},
			"$set":  bson.M{"metadata": message.Metadata},
		}
		result := mongo.FindOneAndUpdate(filter, update, &opt)
		common.Check(result.Err())
		logging(result)
	}
}

func upsertOptions() options.FindOneAndUpdateOptions {
	upsert := true
	after := options.After
	return options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
}

func logging(result *mongo2.SingleResult) {
	var doc Document // := bson.M{}
	var h2h string
	result.Decode(&doc)

	if doc.Metadata.H2H != "" {
		h2h = doc.Metadata.H2H
	} else {
		h2h = "N/A"
	}

	message := fmt.Sprintf("New Message from %s with length %d", h2h, len(doc.Data.Comments))
	log.Println(message)
}
