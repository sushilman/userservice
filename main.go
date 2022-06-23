package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/api/apiuser"
)

const (
	serviceName = "User Service"
	defaultport = "8080"
)

func main() {
	// log level should be set using environment variable
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	fmt.Println("FACEIT User Service")

	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Str("service", serviceName).Logger()
	ctx := context.Background()

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "1",
		})
	})

	apiuser.InitRoutes(&ctx, &logger, router)

	router.Run()
}
