package apiuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/usererrors"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
)

const (
	DefaultLimit = 16
)

func GetUsersHandler(ctx context.Context, logger *zerolog.Logger, userService *services.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryParams models.GetUserQueryParams
		errBind := c.BindQuery(&queryParams)

		if queryParams.Limit == 0 {
			queryParams.Limit = DefaultLimit
		}

		if errBind != nil {
			logger.Error().Err(errBind).Msg("GetUserHandler: Error when binding the query parameters")
			c.JSON(http.StatusBadRequest, usererrors.NewBadRequestErrorResponse("bad query parameters"))
			return
		}

		users, err := userService.GetUsers(ctx, logger, queryParams)
		if err != nil {
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

func GetUserByIdHandler(ctx context.Context, logger *zerolog.Logger, userService *services.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := userService.GetUserById(ctx, logger, c.Param("userId"))
		if err != nil {
			switch err.(type) {
			case *usererrors.NotFoundError:
				c.JSON(http.StatusNotFound, usererrors.NewNotFoundErrorResponse("user not found"))
				return
			}

			c.JSON(http.StatusInternalServerError, usererrors.NewInternalServerError("Something went wrong"))
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
