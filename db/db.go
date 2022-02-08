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
	Database               string
	MatchesCollection      string
	AssetsCollection       string
	UsersCollection        string
	TransactionsCollection string
	Hostname               string
	ReplicaSet             string
	DirectConnection       bool
	Port                   int
}

// Database TODO: slice of collections
type Database struct {
	Client       *mongo.Client
	Matches      *Collection
	Assets       *Collection
	Users        *Collection
	Transactions *Collection
}

func Initialize() {
	//connect()
	//ping()
	//initializeCollections()
}

func DefaultConnectionOpt(configuration *configuration.Configuration) *Options {
	return &Options{
		Port:                   configuration.Mongo.Port,
		Hostname:               configuration.Mongo.Hostname,
		Database:               configuration.Mongo.Database,
		ReplicaSet:             configuration.Mongo.ReplicaSet,
		UsersCollection:        configuration.Mongo.UsersCollection,
		DirectConnection:       configuration.Mongo.DirectConnection,
		AssetsCollection:       configuration.Mongo.AssetsCollection,
		MatchesCollection:      configuration.Mongo.MatchesCollection,
		TransactionsCollection: configuration.Mongo.TransactionsCollection,
	}
}

func Connect(dbOptions *Options) *Database {
	uri := fmt.Sprintf("mongodb://%s:%d/%s", dbOptions.Hostname, dbOptions.Port, dbOptions.Database) // TODO: support replica discovery if master fails
	options := options.Client().
		ApplyURI(uri).
		SetReplicaSet(dbOptions.ReplicaSet).
		SetDirect(dbOptions.DirectConnection)
	client, err := mongo.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}
	//ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	err = client.Connect(context.TODO())
	common.Check(err)
	assets := &Collection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.AssetsCollection),
	}
	matches := &Collection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.MatchesCollection),
	}
	users := &Collection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.UsersCollection),
	}
	transactions := &Collection{
		Collection: client.Database(dbOptions.Database).Collection(dbOptions.TransactionsCollection),
	}

	return &Database{
		Client:       client,
		Assets:       assets,
		Matches:      matches,
		Users:        users,
		Transactions: transactions,
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
