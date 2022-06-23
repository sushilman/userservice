package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/api/apiuser"
	"github.com/sushilman/userservice/db"
)

const (
	serviceName = "User Service"
	defaultport = "8080"
)

func main() {

	// initializing the logger
	// log level should be set using environment variable
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Str("service", serviceName).Logger()

	ctx := context.Background()

	// Establishing database connection and initializing the DB
	dbUri := os.Getenv("DB_URI")
	database := db.InitDB(ctx, &logger, dbUri)

	userStorage := db.NewStorage(database)

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "1",
		})
	})

	apiuser.InitRoutes(ctx, &logger, router, userStorage)

	router.Run()
}
