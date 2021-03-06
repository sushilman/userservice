// Handler for DELETE /v1/users/:userId

package apiuser

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/services"
	"github.com/sushilman/userservice/usererrors"
)

func DeleteUserByIdHandler(userService services.IUserService) func(c *gin.Context) {
	return func(c *gin.Context) {

		err := userService.DeleteUserById(c, c.Param("userId"))
		if err != nil {
			log.Printf("Error when deleting the user. Error: %+v", err)
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
