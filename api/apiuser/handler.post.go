package apiuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sushilman/userservice/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
	"gitlab.valiton.com/cidm/services-commons-api/apierrors"
)

const (
	BASE_PATH = "/v1/users"
)

func PostUserHandler(ctx context.Context, logger *zerolog.Logger, userService *services.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			fmt.Printf("PostUserHandler: Error when binding request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, apierrors.NewBadRequestErrorResponse())
			return
		}

		userId, err := userService.CreateUser(ctx, logger, userCreation)
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		response := models.UserCreationResponse{
			Link: BASE_PATH + "/" + userId,
		}

		c.JSON(http.StatusCreated, response)
	}
}
