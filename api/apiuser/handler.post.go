package apiuser

import (
	"fmt"
	"net/http"

	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/usererrors"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/services"
)

const (
	BASE_PATH = "/v1/users"
)

func PostUserHandler(userService services.IUserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			fmt.Printf("PostUserHandler: Error when binding request body")
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("bad payload"))
			return
		}

		userId, err := userService.CreateUser(userCreation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		response := models.UserCreationResponse{
			Link: BASE_PATH + "/" + userId,
		}

		c.JSON(http.StatusCreated, response)
	}
}
