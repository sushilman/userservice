package db

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	DB_TIMEOUT = 20 * time.Second
)

func InitDB(ctx context.Context, logger *zerolog.Logger, dbUri string) *mongo.Database {
	connectionString, errParse := connstring.ParseAndValidate(dbUri)
	if errParse != nil {
		logger.Fatal().Err(errParse).Msg("Bad DB connection string")
	}

	ctx, cancel := context.WithTimeout(ctx, DB_TIMEOUT)
	defer cancel()

	client, errMongo := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if errMongo != nil {
		logger.Fatal().Err(errMongo).Msg("Could not connect to DB")
	}

	return client.Database(connectionString.Database)
}
