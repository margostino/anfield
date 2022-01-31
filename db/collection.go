package db

import (
	"context"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBCollection struct {
	Collection *mongo.Collection
}

func (c *DBCollection) Upsert(filter interface{}, document interface{}) *domain.Document {
	options := upsertOptions()
	result := c.Collection.FindOneAndUpdate(context.TODO(), filter, document, options)
	common.Check(result.Err())
	return decode(result)
}

func (c *DBCollection) FindOne(filter interface{}) *domain.Document {
	options := findOneOptions()
	result := c.Collection.FindOne(context.TODO(), filter, options)
	//common.Check(result.Err()) // TODO: verify result. This fails in case of different error
	return decode(result)
}

func upsertOptions() *options.FindOneAndUpdateOptions {
	upsert := true
	after := options.After
	return &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
}

func findOneOptions() *options.FindOneOptions {
	//returnKey := true
	return &options.FindOneOptions{
		//ReturnKey: &returnKey,
	}
}

func decode(result *mongo.SingleResult) *domain.Document {
	var document domain.Document
	result.Decode(&document)
	return &document
}
