package db

import (
	"context"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	Collection *mongo.Collection
}

func (c *Collection) Upsert(filter bson.M, update bson.M, document interface{}) error {
	options := upsertOptions(true)
	result := c.Collection.FindOneAndUpdate(context.TODO(), filter, update, options)
	return decode(result, document)
}

func (c *Collection) UpdateWithContext(filter bson.M, update bson.M, document interface{}, context context.Context) error {
	options := upsertOptions(false)
	result := c.Collection.FindOneAndUpdate(context, filter, update, options)
	return decode(result, document)
}

func (c *Collection) Update(filter bson.M, update bson.M, document interface{}) error {
	return c.UpdateWithContext(filter, update, document, context.TODO())
}

func (c *Collection) InsertWithContext(document interface{}, context context.Context) error {
	options := insertOptions()
	_, err := c.Collection.InsertOne(context, document, options)
	return err
}

func (c *Collection) Insert(document interface{}) error {
	return c.InsertWithContext(document, context.TODO())
}

func (c *Collection) FindOne(filter bson.M, document interface{}) error {
	options := findOneOptions()
	result := c.Collection.FindOne(context.TODO(), filter, options)
	return decode(result, document)
}

func upsertOptions(upsert bool) *options.FindOneAndUpdateOptions {
	after := options.After
	return &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
}

func updateOptions() *options.UpdateOptions {
	return &options.UpdateOptions{}
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
func UpsertAssets(asset *domain.Asset) (bson.M, bson.M) {
	filter := FilterBy(asset.Name)
	update := UpdateAssetQuery(asset)
	return filter, update
}

func UpdateWallet(id string, value float64) (bson.M, bson.M) {
	filter := bson.M{"_id": id}
	update := UpdateWalletQuery(value)
	return filter, update
}

func UpsertMatch(match *domain.Match) (bson.M, bson.M) {
	_, _, identifier := common.ExtractTeamsFrom(match.Metadata.Url)
	filter := FilterBy(identifier)
	update := UpdateMatchQuery(match)
	return filter, update
}

func UpsertMatchCompletion(match *domain.Match) (bson.M, bson.M) {
	filter := MatchFilter(match.Metadata.Url)
	update := UpdateCompletionQuery(match)
	return filter, update
}

func MatchFilter(url string) bson.M {
	_, _, identifier := common.ExtractTeamsFrom(url)
	return FilterBy(identifier)
}

func decode(result *mongo.SingleResult, document interface{}) error {
	if common.IsError(result.Err()) {
		return result.Err()
	}
	result.Decode(document)
	return nil
}

func IsDuplicatedWrite(err error) bool {
	writeError := err.(mongo.WriteException)
	return writeError.HasErrorCode(11000)
}