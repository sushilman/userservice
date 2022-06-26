// Handler for POST /v1/users

package apiuser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
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
			log.Printf("Error when binding request body. Error: %+v", errBindJSON)
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("bad payload"))
			return
		}

		invalidFields := validatePayload(userCreation)
		if len(invalidFields) != 0 {
			log.Println("Payload validation error. Invalid fields: ", invalidFields)
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse(fmt.Sprint("Validation Error. Invalid fields: ", invalidFields)))
			return
		}

		userId, err := userService.CreateUser(c, userCreation)
		if err != nil {
			log.Printf("Error when creating the user. Error: %+v", err)
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		response := models.UserCreationResponse{
			Link: BASE_PATH + "/" + *userId,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// Returns a list of invalid fields
func validatePayload(userCreation models.UserCreation) []string {
	validate := validator.New()
	errValidate := validate.Struct(userCreation)
	invalidFieldTags := make([]string, 0)
	if errValidate != nil {
		if _, ok := errValidate.(*validator.InvalidValidationError); ok {
			log.Printf("Error while validating, Error %+v", errValidate)
			return nil
		}

		for _, err := range errValidate.(validator.ValidationErrors) {
			invalidFieldTags = append(invalidFieldTags, err.StructField())
		}
	}
	return invalidFieldTags
}
