// Service layer that separates the api layer and the DB layer
package services

import (
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
	CreateUser(models.UserCreation) (*string, error)
	DeleteUserById(string) error
	GetUsers(models.GetUserQueryParams) ([]models.User, error)
	GetUserById(string) (*models.User, error)
	UpdateUser(string, models.UserCreation) error
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

func (us *UserService) CreateUser(userCreation models.UserCreation) (*string, error) {
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

	err := us.storage.Insert(user)
	if err != nil {
		return nil, err
	}

	us.broker.Publish(events.USER_CREATED_TOPIC, events.UserCreatedEvent(user))

	return &newUserId, nil
}

func (us *UserService) GetUsers(queryParams models.GetUserQueryParams) ([]models.User, error) {
	users, err := us.storage.GetAll(queryParams)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUserById(userId string) (*models.User, error) {
	user, err := us.storage.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdateUser(userId string, userCreation models.UserCreation) error {
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

	err := us.storage.Update(user)

	if err != nil {
		return err
	}

	us.broker.Publish(events.USER_UPDATED_TOPIC, events.UserUpdatedEvent(user))

	return nil
}

func (us *UserService) DeleteUserById(userId string) error {
	err := us.storage.DeleteById(userId)
	if err != nil {
		return err
	}

	us.broker.Publish(events.USER_DELETED_TOPIC, events.UserDeletedEvent{Id: userId})

	return nil
}
