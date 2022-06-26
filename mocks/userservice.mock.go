// Mock for services.UserService
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sushilman/userservice/models"
)

type MockUserService struct {
	*mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user models.UserCreation) (string, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserService) DeleteUserById(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserService) GetUsers(ctx context.Context, q models.GetUserQueryParams) ([]models.User, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userId string, user models.UserCreation) error {
	args := m.Called(ctx, userId, user)
	return args.Error(0)
}
