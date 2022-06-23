package apiuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
	"gitlab.valiton.com/cidm/services-commons-api/apierrors"
)

func UpdateUserByIdHandler(ctx context.Context, logger *zerolog.Logger, s db.IStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			fmt.Printf("PostUserHandler: Error when binding request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, apierrors.NewBadRequestErrorResponse())
			return
		}

		userService := services.NewUserService(s)

		err := userService.UpdateUserById(logger, c.Param("userId"), userCreation)
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		c.Status(http.StatusNoContent)
	}
}
