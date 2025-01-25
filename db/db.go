package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI      string
	Database string
}

func InitializeMongoDB(config MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(config.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Successfully connected to MongoDB!")

	database := client.Database(config.Database)

	return client, database, nil
}
