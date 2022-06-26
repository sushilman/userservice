// Handler for GET /v1/users and GET /v1/users/:userId

package apiuser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/usererrors"

	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/services"
)

const (
	DEFAULT_PAGE_LIMIT = 16
)

func GetUsersHandler(userService services.IUserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryParams models.GetUserQueryParams
		errBind := c.BindQuery(&queryParams)

		if queryParams.Limit == 0 {
			queryParams.Limit = DEFAULT_PAGE_LIMIT
		}

		if errBind != nil {
			log.Printf("Error when binding the query parameters.\nError: %+v", errBind)
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("bad query parameters"))
			return
		}

		users, err := userService.GetUsers(c, queryParams)
		if err != nil {
			log.Printf("Error while fetching users.\nError: %+v", err)
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		prevOffset := queryParams.Offset - queryParams.Limit
		if queryParams.Limit > queryParams.Offset {
			prevOffset = 0
		}

		prevPage := fmt.Sprintf("%s?offset=%d&limit=%d", BASE_PATH, prevOffset, queryParams.Limit)
		selfPage := fmt.Sprintf("%s?offset=%d&limit=%d", BASE_PATH, queryParams.Offset, queryParams.Limit)
		nextPage := fmt.Sprintf("%s?offset=%d&limit=%d", BASE_PATH, queryParams.Offset+queryParams.Limit, queryParams.Limit)

		if queryParams.Offset == 0 {
			prevPage = ""
		}

		// TODO: handle when the returned users are exactly the `limit` size and it is the last page
		// the idea could be to always fetch `limit+1` from DB but return only up to `limit` in the reponse
		if len(users) < int(queryParams.Limit) {
			nextPage = ""
		}

		response := models.GetUsersResponse{
			Data: users,
			Links: models.PaginationLinks{
				Prev: prevPage,
				Self: selfPage,
				Next: nextPage,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetUserByIdHandler(userService services.IUserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := userService.GetUserById(c, c.Param("userId"))
		if err != nil {
			switch err.(type) {
			case *usererrors.NotFoundError:
				c.JSON(http.StatusNotFound, usererrors.NewNotFoundErrorResponse("user not found"))
				return
			}

			log.Printf("Error while fetching user by ID.\nError: %+v", err)
			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
