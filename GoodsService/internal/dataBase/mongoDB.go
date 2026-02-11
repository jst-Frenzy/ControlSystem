package dataBase

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var MongoDB *mongo.Client

func InitMongo() {
	var err error

	uri := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(uri)

	opts.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	MongoDB, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal("Failed to connect to Mongo")
	}
}
