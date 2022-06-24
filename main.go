package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/api/apiuser"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/messagebroker"
	"github.com/sushilman/userservice/services"
)

const (
	SERVICE_NAME = "User Service"
	DEFAULT_PORT = "8080"
)

func main() {
	// Establishing database connection and initializing the DB
	dbUri := os.Getenv("DB_URI")
	database := db.InitDB(dbUri)

	// Initializing message broker
	messagebroker := messagebroker.InitMessageBroker()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "1",
		})
	})

	userStorage := db.NewStorage(database)
	userservice := services.NewUserService(userStorage, messagebroker)
	apiuser.InitRoutes(router, userservice)

	router.Run()
}
