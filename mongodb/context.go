package mongo

import (
	"context"
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var mongoConnection *mongo.Client
var mongoCollection *mongo.Collection
var ctx context.Context

func Initialize() {
	//var mongoDatabase *mongo.Database
	var err error

	hostname := configuration.Mongo().Hostname
	database := configuration.Mongo().Database
	matchesCollection := configuration.Mongo().MatchesCollection
	port := configuration.Mongo().Port
	uri := fmt.Sprintf("mongodb://%s:%d", hostname, port)
	mongoConnection, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoConnection.Connect(ctx)
	common.Check(err)
	err = mongoConnection.Ping(ctx, readpref.Primary())
	common.Check(err)
	//databases, err := mongoConnection.ListDatabaseNames(ctx, bson.M{})
	//common.Check(err)
	//fmt.Println(databases)

	mongoCollection = mongoConnection.Database(database).Collection(matchesCollection)
}

func Context() context.Context {
	return ctx
}

func Close() {
	mongoConnection.Disconnect(ctx)
}

func Insert(document interface{}) (*mongo.InsertOneResult, error) {
	return mongoCollection.InsertOne(ctx, document)
}

func FindOneAndUpdate(filter interface{}, document interface{}, options *options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return mongoCollection.FindOneAndUpdate(ctx, filter, document, options)
}

func Find(filter interface{}) (*mongo.Cursor, error) {
	return mongoCollection.Find(ctx, filter)
}
