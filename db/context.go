package db

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

type Options struct {
	Database          string
	MatchesCollection string
	ScoresCollection  string
	Hostname          string
	Port              int
}

type Database struct {
	Client  *mongo.Client
	Matches *DBCollection
	Scores  *DBCollection
}

func Initialize() {
	//connect()
	//ping()
	//initializeCollections()
}

func DefaultConnectionOpt(configuration *configuration.Configuration) *Options {
	return &Options{
		Database:          configuration.Mongo.Database,
		MatchesCollection: configuration.Mongo.MatchesCollection,
		ScoresCollection:  configuration.Mongo.ScoresCollection,
		Hostname:          configuration.Mongo.Hostname,
		Port:              configuration.Mongo.Port,
	}
}

func Connect(dbOptions *Options) *Database {
	uri := fmt.Sprintf("mongodb://%s:%d", dbOptions.Hostname, dbOptions.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	//ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	err = client.Connect(context.TODO())
	common.Check(err)
	scores := &DBCollection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.ScoresCollection),
	}
	matches := &DBCollection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.MatchesCollection),
	}

	return &Database{
		Client:  client,
		Scores:  scores,
		Matches: matches,
	}
}

func (d *Database) ping() {
	err := d.Client.Ping(context.TODO(), readpref.Primary())
	common.Check(err)
	//databases, err := mongoConnection.ListDatabaseNames(ctx, bson.M{})
	//common.Check(err)
	//fmt.Println(databases)
}

func (d *Database) Close() {
	d.Client.Disconnect(context.TODO())
}
