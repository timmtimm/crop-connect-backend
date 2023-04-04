package mongo_driver

import (
	"context"
	"fmt"
	"marketplace-backend/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(databaseName string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Make MongoDB URI
	dbUser := util.GetConfig("DB_USER")
	dbPass := util.GetConfig("DB_PASS")
	dbHost := util.GetConfig("DB_HOST")
	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", dbUser, dbPass, dbHost)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongodbURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	database := client.Database(databaseName)
	return database
}

func Close(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return db.Client().Disconnect(ctx)
}

func CheckCollectionExist(db *mongo.Database, collectionName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return false, err
	}

	fmt.Println()

	return len(cursor) > 0, nil
}
