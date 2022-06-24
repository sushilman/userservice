package apiuser

import (
	"github.com/gin-gonic/gin"
	"github.com/sushilman/userservice/services"
)

func InitRoutes(r *gin.Engine, us services.IUserService) {
	r.POST("/v1/users", PostUserHandler(us))
	r.GET("/v1/users", GetUsersHandler(us))
	r.GET("/v1/users/:userId", GetUserByIdHandler(us))
	r.PUT("/v1/users/:userId", UpdateUserHandler(us))
	r.DELETE("/v1/users/:userId", DeleteUserByIdHandler(us))
}
