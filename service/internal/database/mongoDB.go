package database

import (
	"News/service/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func newMongoDB(ctx context.Context, dbName string, db config.MongoDB) *mongo.Database {
	uri := fmt.Sprintf("mongodb://%s:%s/%s", db.Host, db.Port, db.Database)
	clientOptions := options.Client()
	clientOptions.ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %s, err: %v", dbName, err)
	}
	log.Printf("Connected to MongoDB: %s", dbName)

	return client.Database(db.Database)
}
