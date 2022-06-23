package services_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
	// "github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	userCreation := models.UserCreation{}
	testLogger := zerolog.Nop()

	userservice := services.NewUserService()

	userId, err := userservice.CreateUser(&testLogger, userCreation)

	require.Nil(t, err)
	assert.NotNil(t, userId)
}

func TestGetUser(t *testing.T) {
	// TODO: mock tests
}

func TestUpdateUser(t *testing.T) {
	// TODO:
}

func TestDeleteUser(t *testing.T) {
	// TODO:
}
