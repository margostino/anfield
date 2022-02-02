package db

import (
	"context"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Collection struct {
	Collection *mongo.Collection
}

func (c *Collection) UpsertMatch(filter interface{}, document interface{}) *domain.MatchDocument {
	result := c.upsert(filter, document)
	common.Check(result.Err())
	return decodeMatch(result)
}

func (c *Collection) UpsertAsset(filter interface{}, document interface{}) *domain.AssetDocument {
	result := c.upsert(filter, document)
	common.Check(result.Err())
	return decodeAsset(result)
}

func (c *Collection) UpsertUser(filter interface{}, document interface{}) *domain.UserDocument {
	result := c.upsert(filter, document)
	common.Check(result.Err())
	return decodeUser(result)
}

func (c *Collection) InsertUser(document interface{}) error {
	_, err := c.insert(document)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (c *Collection) FindOneMatch(filter interface{}) *domain.MatchDocument {
	options := findOneOptions()
	result := c.Collection.FindOne(context.TODO(), filter, options)
	//common.Check(result.Err()) // TODO: verify result. This fails in case of different error
	return decodeMatch(result)
}

func (c *Collection) FindOneAsset(filter interface{}) *domain.AssetDocument {
	options := findOneOptions()
	result := c.Collection.FindOne(context.TODO(), filter, options)
	//common.Check(result.Err()) // TODO: verify result. This fails in case of different error
	return decodeAsset(result)
}

func (c *Collection) upsert(filter interface{}, document interface{}) *mongo.SingleResult {
	options := upsertOptions()
	result := c.Collection.FindOneAndUpdate(context.TODO(), filter, document, options)
	return result
}

func (c *Collection) insert(document interface{}) (*mongo.InsertOneResult, error) {
	options := insertOptions()
	return c.Collection.InsertOne(context.TODO(), document, options)
}

func upsertOptions() *options.FindOneAndUpdateOptions {
	upsert := true
	after := options.After
	return &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
}

func insertOptions() *options.InsertOneOptions {
	return &options.InsertOneOptions{}
}

func findOneOptions() *options.FindOneOptions {
	//returnKey := true
	return &options.FindOneOptions{
		//ReturnKey: &returnKey,
	}
}

// TODO: generics to reduce boilerplate

func decodeMatch(result *mongo.SingleResult) *domain.MatchDocument {
	var document domain.MatchDocument
	result.Decode(&document)
	return &document
}

func decodeAsset(result *mongo.SingleResult) *domain.AssetDocument {
	var document domain.AssetDocument
	result.Decode(&document)
	return &document
}

func decodeUser(result *mongo.SingleResult) *domain.UserDocument {
	var document domain.UserDocument
	result.Decode(&document)
	return &document
}

func isDuplicatedWrite(err error) bool {
	writeError := err.(mongo.WriteException)
	return writeError.HasErrorCode(11000)
}
