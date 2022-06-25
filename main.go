package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/api/apiuser"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/messagebroker"
	"github.com/sushilman/userservice/services"
)

const (
	SERVICE_NAME              = "User Service"
	DEFAULT_PORT              = "8080"
	GRACEFUL_SHUTDOWN_TIMEOUT = 25
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

	srv := &http.Server{
		Addr:    ":" + DEFAULT_PORT,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Unable to start the server")
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(GRACEFUL_SHUTDOWN_TIMEOUT)*time.Second)
	defer cancel()
	db.CloseDB(ctx, database)
}
