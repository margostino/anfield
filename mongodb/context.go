package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var mongoConnection *mongo.Client
var ctx context.Context

func Initialize() {
	var err error
	mongoConnection, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoConnection.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	mongoConnection.Disconnect(ctx)
}
