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

var client *mongo.Client
var Matches, Scores *DBCollection
var ctx context.Context

func Initialize() {
	startClient()
	connect()
	ping()
	initializeCollections()
}

func ping() {
	err := client.Ping(ctx, readpref.Primary())
	common.Check(err)
	//databases, err := mongoConnection.ListDatabaseNames(ctx, bson.M{})
	//common.Check(err)
	//fmt.Println(databases)
}

func connect() {
	//ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	ctx = context.TODO()
	err := client.Connect(ctx)
	common.Check(err)
}

func startClient() {
	uri := uri()
	cli, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	client = cli
}

func uri() string {
	hostname := configuration.Mongo().Hostname
	port := configuration.Mongo().Port
	return fmt.Sprintf("mongodb://%s:%d", hostname, port)
}

func initializeCollections() {
	database := configuration.Mongo().Database
	matchesCollection := configuration.Mongo().MatchesCollection
	scoresCollection := configuration.Mongo().ScoresCollection
	Scores = &DBCollection{
		Collection: client.Database(database).Collection(scoresCollection),
	}
	Matches = &DBCollection{
		Collection: client.Database(database).Collection(matchesCollection),
	}
}

func Close() {
	client.Disconnect(ctx)
}
