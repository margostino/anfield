package dataloader

import (
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/kafka"
	mongo "github.com/margostino/anfield/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
