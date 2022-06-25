// Handle the Mongo DB connection

package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	DB_TIMEOUT = 20 * time.Second
)

func InitDB(dbUri string) *mongo.Database {
	connectionString, errParse := connstring.ParseAndValidate(dbUri)
	if errParse != nil {
		fmt.Printf("Bad DB connection string: %v", dbUri)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	client, errMongo := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if errMongo != nil {
		fmt.Printf("Could not connect to DB: %v", dbUri)
		os.Exit(1)
	}

	return client.Database(connectionString.Database)
}

func CloseDB(ctx context.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(ctx, DB_TIMEOUT)
	defer cancel()

	if err := db.Client().Disconnect(ctx); err != nil {
		fmt.Printf("Error while closing DB connection. Error: %+v", err)
	}
}
