package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var Client *mongo.Client

func DbConnection() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return CreateIndexes()
}

func CreateIndexes() error {
	if Client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"userMessage": 1},
		Options: options.Index().SetUnique(false), 
	}

	historyCollection := Client.Database("Chat").Collection("history")

	_, err := historyCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index on 'history' collection: %v", err)
	}

	fmt.Println("Index created successfully on 'history' collection")
	return nil
}
