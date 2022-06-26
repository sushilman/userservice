package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/api/apiuser"
	"github.com/sushilman/userservice/db"
	pb "github.com/sushilman/userservice/grpc/proto"
	userservice "github.com/sushilman/userservice/grpc/server"
	"github.com/sushilman/userservice/messagebroker"
	"github.com/sushilman/userservice/services"
	"google.golang.org/grpc"
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

	// Start http server
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Unable to start the server")
		}
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	listen, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Printf("Failed to listen")
	}

	pb.RegisterUserServiceServer(grpcServer, newgRPCUserServer(userservice))
	grpcServer.Serve(listen)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(GRACEFUL_SHUTDOWN_TIMEOUT)*time.Second)
	defer cancel()
	db.CloseDB(ctx, database)
}

func newgRPCUserServer(svc services.IUserService) *userservice.UserServicegRPCServer {
	return &userservice.UserServicegRPCServer{
		US: svc,
	}
}
