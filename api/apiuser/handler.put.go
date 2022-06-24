package apiuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
	"github.com/sushilman/userservice/usererrors"
)

func UpdateUserHandler(ctx context.Context, logger *zerolog.Logger, userService *services.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			fmt.Printf("PostUserHandler: Error when binding request body")
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("bad payload"))
			return
		}

		err := userService.UpdateUser(ctx, logger, c.Param("userId"), userCreation)
		if err != nil {
			switch err.(type) {
			case *usererrors.NotFoundError:
				c.JSON(http.StatusNotFound, usererrors.NewNotFoundErrorResponse("User not found"))
				return
			}

			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
