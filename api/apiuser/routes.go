package apiuser

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitRoutes(ctx *context.Context, logger *zerolog.Logger, e *gin.Engine) {
	e.POST("/v1/users", PostUserHandler(ctx, logger))
}
