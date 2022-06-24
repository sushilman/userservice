package apiuser

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
	"github.com/sushilman/userservice/usererrors"
)

func DeleteUserByIdHandler(ctx context.Context, logger *zerolog.Logger, userService *services.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {

		err := userService.DeleteUserById(ctx, logger, c.Param("userId"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
