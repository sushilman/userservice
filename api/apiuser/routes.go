package apiuser

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitRoutes(logger *zerolog.Logger, e *gin.Engine) {
	e.POST("/v1/users", PostUserHandler(logger))
	e.GET("/v1/users", GetUsersHandler(logger))
	e.GET("/v1/users/:userId", GetUserByIdHandler(logger))
	e.PUT("/v1/users/:userId", UpdateUserByIdHandler(logger))
}
