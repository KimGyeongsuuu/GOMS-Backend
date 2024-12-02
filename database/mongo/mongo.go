package mongo

import (
	"GOMS-BACKEND-GO/global/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection() (*mongo.Client, *mongo.Database, error) {
	user := config.Data().Mongo.User
	password := config.Data().Mongo.Pass
	host := config.Data().Mongo.Host
	port := config.Data().Mongo.Port
	database := config.Data().Mongo.Db

	var uri string
	if user != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", user, password, host, port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%d", host, port)
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, client.Database(database), nil
}

func CreateCollections(db *mongo.Database, collectionNames []string) {
	for _, name := range collectionNames {
		collection := db.Collection(name)

		_, err := collection.InsertOne(context.Background(), bson.D{})
		if err != nil {
			log.Printf("Failed to create collection %s: %v", name, err)
		} else {
			fmt.Printf("Created collection: %s\n", name)
		}
	}
}
