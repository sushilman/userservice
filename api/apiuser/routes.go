package apiuser

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/db"
)

func InitRoutes(ctx context.Context, logger *zerolog.Logger, r *gin.Engine, s db.IStorage) {
	r.POST("/v1/users", PostUserHandler(ctx, logger, s))
	r.GET("/v1/users", GetUsersHandler(ctx, logger, s))
	r.GET("/v1/users/:userId", GetUserByIdHandler(ctx, logger, s))
	r.PUT("/v1/users/:userId", UpdateUserHandler(ctx, logger, s))
	r.DELETE("/v1/users/:userId", DeleteUserByIdHandler(ctx, logger, s))
}
