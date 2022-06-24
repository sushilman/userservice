package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/events"
	"github.com/sushilman/userservice/messagebroker"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/utils"
)

type UserService struct {
	storage db.IUserStorage
	broker  messagebroker.IMessageBroker
}

func NewUserService(storage db.IUserStorage, broker messagebroker.IMessageBroker) *UserService {
	return &UserService{storage, broker}
}

func (us *UserService) CreateUser(ctx context.Context, logger *zerolog.Logger, userCreation models.UserCreation) (string, error) {
	newUserId := uuid.NewString()
	createdAt := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, errHash := utils.HashPassword(userCreation.Password)
	if errHash != nil {
		logger.Error().Err(errHash).Msg("Error while hashing password")
		return "", errHash
	}

	user := models.User{
		Id:        newUserId,
		FirstName: userCreation.FirstName,
		LastName:  userCreation.LastName,
		Nickname:  userCreation.Nickname,
		Password:  hashedPassword,
		Email:     userCreation.Email,
		Country:   userCreation.Country,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	err := us.storage.Insert(ctx, user)
	if err != nil {
		return "", err
	}

	us.broker.Publish(logger, events.USER_CREATED_TOPIC, events.UserCreatedEvent(user))

	return newUserId, nil
}

func (us *UserService) GetUsers(ctx context.Context, logger *zerolog.Logger, queryParams models.GetUserQueryParams) ([]models.User, error) {
	users, err := us.storage.GetAll(ctx, queryParams)
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

func (us *UserService) UpdateUser(ctx context.Context, logger *zerolog.Logger, userId string, userCreation models.UserCreation) error {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, errHash := utils.HashPassword(userCreation.Password)
	if errHash != nil {
		logger.Error().Err(errHash).Msg("Error while hashing password")
		return errHash
	}

	user := models.User{
		Id:        userId,
		FirstName: userCreation.FirstName,
		LastName:  userCreation.LastName,
		Nickname:  userCreation.Nickname,
		Password:  hashedPassword,
		Email:     userCreation.Email,
		Country:   userCreation.Country,
		UpdatedAt: updatedAt,
	}

	err := us.storage.Update(ctx, user)

	if err != nil {
		return err
	}

	us.broker.Publish(logger, events.USER_UPDATED_TOPIC, events.UserUpdatedEvent(user))

	return nil
}

func (us *UserService) DeleteUserById(ctx context.Context, logger *zerolog.Logger, userId string) error {
	err := us.storage.DeleteById(ctx, userId)
	if err != nil {
		return err
	}

	us.broker.Publish(logger, events.USER_UPDATED_TOPIC, events.UserDeletedEvent{Id: userId})

	return nil
}
