// Service layer that separates the api layer and the DB layer
package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sushilman/userservice/db"
	"github.com/sushilman/userservice/events"
	"github.com/sushilman/userservice/messagebroker"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/utils"
)

type IUserService interface {
	CreateUser(context.Context, models.UserCreation) (*string, error)
	DeleteUserById(context.Context, string) error
	GetUsers(context.Context, models.GetUserQueryParams) ([]models.User, error)
	GetUserById(context.Context, string) (*models.User, error)
	UpdateUser(context.Context, string, models.UserCreation) error
}

type UserService struct {
	storage db.IUserStorage
	broker  messagebroker.IMessageBroker
}

func NewUserService(storage db.IUserStorage, broker messagebroker.IMessageBroker) IUserService {
	return &UserService{
		storage,
		broker,
	}
}

// CreateUser returns the ID of the newly created user
func (us *UserService) CreateUser(ctx context.Context, userCreation models.UserCreation) (*string, error) {
	newUserId := uuid.NewString()
	createdAt := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, errHash := utils.HashPassword(userCreation.Password)
	if errHash != nil {
		log.Printf("Error while hashing password. Error: %+v", errHash)
		return nil, errHash
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
		return nil, err
	}

	us.broker.Publish(events.USER_CREATED_TOPIC, events.UserCreatedEvent(user))

	return &newUserId, nil
}

// GetUsers returns the list of the Users, filtered by parameters in the queryParams
func (us *UserService) GetUsers(ctx context.Context, queryParams models.GetUserQueryParams) ([]models.User, error) {
	users, err := us.storage.GetAll(ctx, queryParams)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserById returns the user having the provided userId
func (us *UserService) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	user, err := us.storage.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates the user having the provided userId
func (us *UserService) UpdateUser(ctx context.Context, userId string, userCreation models.UserCreation) error {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, errHash := utils.HashPassword(userCreation.Password)
	if errHash != nil {
		log.Printf("Error while hashing password. Error: %+v", errHash)
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

	us.broker.Publish(events.USER_UPDATED_TOPIC, events.UserUpdatedEvent(user))

	return nil
}

// DeleteUserById deletes the user having the provided userId
func (us *UserService) DeleteUserById(ctx context.Context, userId string) error {
	err := us.storage.DeleteById(ctx, userId)
	if err != nil {
		return err
	}

	us.broker.Publish(events.USER_DELETED_TOPIC, events.UserDeletedEvent{Id: userId})

	return nil
}
