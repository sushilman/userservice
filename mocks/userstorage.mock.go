package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sushilman/userservice/models"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Insert(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStorage) GetAll(queryParams models.GetUserQueryParams) (users []models.User, err error) {
	args := m.Called(queryParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserStorage) GetById(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserStorage) Update(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStorage) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
