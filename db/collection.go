package db

import (
	"context"
	"github.com/margostino/anfield/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBCollection struct {
	Collection *mongo.Collection
}

func (c *DBCollection) Upsert(filter interface{}, document interface{}) *mongo.SingleResult {
	options := upsertOptions()
	result := c.Collection.FindOneAndUpdate(context.TODO(), filter, document, options)
	common.Check(result.Err())
	return result
}

func upsertOptions() *options.FindOneAndUpdateOptions {
	upsert := true
	after := options.After
	return &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
}
