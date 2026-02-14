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

	MongoDB, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal("Failed to connect to Mongo")
	}
}
