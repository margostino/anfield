package dataloader

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

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
		log.Println("Message consumed and stored", message.Metadata.H2H)
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
