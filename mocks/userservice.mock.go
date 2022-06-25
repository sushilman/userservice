// Mock for services.UserService
package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sushilman/userservice/models"
)

type MockUserService struct {
	*mock.Mock
}

func (m *MockUserService) CreateUser(user models.UserCreation) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserService) DeleteUserById(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserService) GetUsers(q models.GetUserQueryParams) ([]models.User, error) {
	args := m.Called(q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserById(userId string) (*models.User, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(userId string, user models.UserCreation) error {
	args := m.Called(userId, user)
	return args.Error(0)
}
