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

func (us *UserService) GetUsers(logger *zerolog.Logger, queryParams models.GetUserQueryParams) ([]models.User, error) {
	// TODO: implement
	// this is just the mocked mockedUsers
	mockedUsers := []models.User{
		{
			Id:        "c499e35b-2586-45a0-a304-ea602f3e9922",
			FirstName: "Maximilian",
			LastName:  "Mustermann",
			Nickname:  "Max",
			Country:   "DE",
			Email:     "max@example.com",
			CreatedAt: "2019-10-12T07:20:50.52Z",
			UpdatedAt: "2019-10-12T07:20:50.52Z",
		},
		{
			Id:        "f960ff54-0de7-4022-9528-7d6e50fe2d7e",
			FirstName: "Samuel",
			LastName:  "Wilson",
			Nickname:  "Sam",
			Country:   "UK",
			Email:     "sam@example.com",
			CreatedAt: "2019-10-12T07:20:50.52Z",
			UpdatedAt: "2019-10-12T07:20:50.52Z",
		},
	}

	return mockedUsers, nil
}

func (us *UserService) GetUserById(logger *zerolog.Logger, userId string) (models.User, error) {
	// TODO: implement

	mockedUser := models.User{
		Id:        userId,
		FirstName: "Maximilian",
		LastName:  "Mustermann",
		Nickname:  "Max",
		Country:   "DE",
		Email:     "max@example.com",
		CreatedAt: "2019-10-12T07:20:50.52Z",
		UpdatedAt: "2019-10-12T07:20:50.52Z",
	}

	return mockedUser, nil
}

func (us *UserService) UpdateUserById(logger *zerolog.Logger, userId string) error {
	// TODO: implement
	return nil
}

func (us *UserService) DeleteUserById(logger *zerolog.Logger, userId string) error {
	// TODO: implement
	return nil
}
