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
	BasePath = "/v1/users"
)

func PostUserHandler(ctx *context.Context, logger *zerolog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			fmt.Printf("PostUserHandler: Error when binding request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, apierrors.NewBadRequestErrorResponse())
			return
		}

		userservice := services.NewUserService()

		userId, err := userservice.CreateUser(ctx, logger, userCreation)
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		response := models.UserCreationResponse{
			Link: BasePath + "/" + userId,
		}

		c.JSON(http.StatusCreated, response)
	}
}
