package apiuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/services"
	"gitlab.valiton.com/cidm/services-commons-api/apierrors"
)

const (
	DefaultLimit = 16
)

func GetUsersHandler(ctx context.Context, logger *zerolog.Logger, s db.IStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryParams models.GetUserQueryParams
		errBind := c.BindQuery(&queryParams)

		if queryParams.Limit == 0 {
			queryParams.Limit = DefaultLimit
		}

		if errBind != nil {
			logger.Error().Err(errBind).Msg("GetUserHandler: Error when binding the query parameters")
			c.AbortWithStatusJSON(http.StatusBadRequest, apierrors.NewBadRequestErrorResponse())
			return
		}

		userService := services.NewUserService(s)

		users, err := userService.GetUsers(ctx, logger, queryParams)
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
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

		// TODO: handle when the returned users are exactly 16 and it is the last page
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

func GetUserByIdHandler(ctx context.Context, logger *zerolog.Logger, s db.IStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		userservice := services.NewUserService(s)

		user, err := userservice.GetUserById(ctx, logger, c.Param("userId"))
		if err != nil {
			logger.Err(err).Msg("Something went wrong. TODO: Handle error")
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
