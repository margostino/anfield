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
)

var mongoConnection *mongo.Client
var matches, scores *mongo.Collection
var ctx context.Context

func Initialize() {
	//var mongoDatabase *mongo.Database
	var err error

	hostname := configuration.Mongo().Hostname

	port := configuration.Mongo().Port
	uri := fmt.Sprintf("mongodb://%s:%d", hostname, port)
	mongoConnection, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	//ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	ctx = context.TODO()
	err = mongoConnection.Connect(ctx)
	common.Check(err)
	err = mongoConnection.Ping(ctx, readpref.Primary())
	common.Check(err)
	//databases, err := mongoConnection.ListDatabaseNames(ctx, bson.M{})
	//common.Check(err)
	//fmt.Println(databases)

	matches, scores = getCollections()

}

func getCollections() (*mongo.Collection, *mongo.Collection){
	database := configuration.Mongo().Database
	matchesCollection := configuration.Mongo().MatchesCollection
	scoresCollection := configuration.Mongo().ScoresCollection
	matches = mongoConnection.Database(database).Collection(matchesCollection)
	scores = mongoConnection.Database(database).Collection(scoresCollection)
	return matches, scores
}

func Close() {
	mongoConnection.Disconnect(ctx)
}

func FindOneAndUpdate(filter interface{}, document interface{}, options *options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return mongoCollection.FindOneAndUpdate(ctx, filter, document, options)
}

func Find(filter interface{}) (*mongo.Cursor, error) {
	return mongoCollection.Find(ctx, filter)
}
