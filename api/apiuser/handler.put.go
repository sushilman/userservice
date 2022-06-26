// Handler for PUT /v1/users/:userId

package apiuser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
	"github.com/sushilman/userservice/usererrors"
)

func UpdateUserHandler(userService services.IUserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userCreation models.UserCreation
		errBindJSON := c.BindJSON(&userCreation)
		if errBindJSON != nil {
			log.Printf("Error when binding request body")
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("Bad Payload"))
			return
		}

		invalidFields := validatePayload(userCreation)
		if len(invalidFields) != 0 {
			log.Println("Payload validation error. Invalid fields: ", invalidFields)
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse(fmt.Sprint("Validation Error. Invalid fields: ", invalidFields)))
			return
		}

		err := userService.UpdateUser(c.Param("userId"), userCreation)
		if err != nil {
			switch err.(type) {
			case *usererrors.NotFoundError:
				c.JSON(http.StatusNotFound, usererrors.NewNotFoundErrorResponse("User not found"))
				return
			}

			log.Printf("Error when updating the user. Error: %+v", err)
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
