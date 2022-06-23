package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/models"
)

type UserService struct {
	storage db.IStorage
}

func NewUserService(storage db.IStorage) *UserService {
	return &UserService{storage}
}

func (us *UserService) CreateUser(ctx context.Context, logger *zerolog.Logger, userCreation models.UserCreation) (string, error) {
	newUserId := uuid.NewString()
	createdAt := time.Now().UTC().Format(time.RFC3339)

	us.storage.Insert(ctx, models.User{
		Id:        newUserId,
		FirstName: userCreation.FirstName,
		LastName:  userCreation.LastName,
		Nickname:  userCreation.Nickname,
		Email:     userCreation.Email,
		Country:   userCreation.Country,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})

	return newUserId, nil
}

func (us *UserService) GetUsers(ctx context.Context, logger *zerolog.Logger, queryParams models.GetUserQueryParams) ([]models.User, error) {
	users, err := us.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUserById(ctx context.Context, logger *zerolog.Logger, userId string) (*models.User, error) {
	user, err := us.storage.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdateUserById(logger *zerolog.Logger, userId string, userCreation models.UserCreation) error {
	// TODO: implement
	return nil
}

func (us *UserService) DeleteUserById(logger *zerolog.Logger, userId string) error {
	// TODO: implement
	return nil
}
