package db

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sushilman/userservice/models"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Insert(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
