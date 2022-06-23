package apiuser

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/services"
)

func DeleteUserByIdHandler(ctx context.Context, logger *zerolog.Logger, s db.IStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		userService := services.NewUserService(s)

		err := userService.DeleteUserById(logger, c.Param("userId"))
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		c.Status(http.StatusNoContent)
	}
}
