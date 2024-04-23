package db

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func DbConnection() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    var err error
    Client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint
        if err := Client.Disconnect(ctx); err != nil {
            log.Printf("Failed to disconnect MongoDB client: %v", err)
        }
        cancel()
    }()

    err = Client.Ping(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    fmt.Println("Connected to MongoDB!")

    return InitializeCollection(ctx)
}

func InitializeCollection(ctx context.Context) error {
    if Client == nil {
        return fmt.Errorf("MongoDB client is not initialized")
    }

    historyCollection := Client.Database("Chat").Collection("history")

    // Check if the collection exists and has documents
    count, err := historyCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return fmt.Errorf("failed to count documents in 'history' collection: %v", err)
    }

    // If documents are found, drop the collection
    if count > 0 {
        err := historyCollection.Drop(ctx)
        if err != nil {
            return fmt.Errorf("failed to drop 'history' collection: %v", err)
        }
        fmt.Println("Dropped existing 'history' collection")
    }

    return CreateIndexes(ctx, historyCollection)
}

func CreateIndexes(ctx context.Context, historyCollection *mongo.Collection) error {
    indexModel := mongo.IndexModel{
        Keys:    bson.M{"userMessage": 1}, // Example index field
        Options: options.Index().SetUnique(false),
    }

    _, err := historyCollection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return fmt.Errorf("failed to create index on 'history' collection: %v", err)
    }

    fmt.Println("Index created successfully on 'history' collection")
    return nil
}