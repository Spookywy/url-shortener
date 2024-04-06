package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDbClient() *mongo.Client {
	uri, ok := os.LookupEnv("DB_URI")
	if !ok {
		log.Fatal("DB_URI environment variable not found")
	}

	timeout := time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Insure that the connection is established
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	return client
}
