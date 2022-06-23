package apiuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
)

func DeleteUserByIdHandler(logger *zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		userservice := services.NewUserService()

		err := userservice.DeleteUserById(logger, c.Param("userId"))
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		c.Status(http.StatusNoContent)
	}
}
