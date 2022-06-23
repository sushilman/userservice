package services_test

import (
	"context"
	"testing"
	// "github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// userCreation := models.UserCreation{}
	// testLogger := zerolog.Nop()
	// testCtx := testContext()

	// mockUserStorage := new(db.MockStorage)
	// mockUserStorage.On("Insert", testCtx, mock.AnythingOfType("models.User")).Return(nil).Times(1)

	// userservice := services.NewUserService(mockUserStorage)
	// userId, err := userservice.CreateUser(testCtx, &testLogger, userCreation)

	// require.Nil(t, err)
	// assert.NotNil(t, userId)
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

func testContext() context.Context {
	return context.Background()
}
