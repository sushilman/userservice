package db

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sushilman/userservice/models"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Insert(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
