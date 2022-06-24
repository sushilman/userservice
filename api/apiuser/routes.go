package apiuser

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
)

func InitRoutes(ctx context.Context, logger *zerolog.Logger, r *gin.Engine, us *services.UserService) {
	r.POST("/v1/users", PostUserHandler(ctx, logger, us))
	r.GET("/v1/users", GetUsersHandler(ctx, logger, us))
	r.GET("/v1/users/:userId", GetUserByIdHandler(ctx, logger, us))
	r.PUT("/v1/users/:userId", UpdateUserHandler(ctx, logger, us))
	r.DELETE("/v1/users/:userId", DeleteUserByIdHandler(ctx, logger, us))
}
