package services

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) CreateUser(logger *zerolog.Logger, userCreation models.UserCreation) (string, error) {
	newUserId := uuid.NewString()
	return newUserId, nil
}
