package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/sushilman/userservice/events"
	"github.com/sushilman/userservice/mocks"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
)

// Test for User Creation
func TestCreateUser(t *testing.T) {
	userCreation := testDataUserCreation()

	expectedUser := models.User{
		FirstName: userCreation.FirstName,
		LastName:  userCreation.LastName,
		Nickname:  userCreation.Nickname,
		Email:     userCreation.Email,
		Country:   userCreation.Country,
	}

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("Insert", mock.MatchedBy(
		func(actualUser models.User) bool {
			expectedUser.Id = actualUser.Id
			expectedUser.Password = actualUser.Password
			expectedUser.CreatedAt = actualUser.CreatedAt
			expectedUser.UpdatedAt = actualUser.UpdatedAt
			return assert.Equal(t, expectedUser, actualUser)
		},
	)).Return(nil)

	expectedUserEvent := events.UserCreatedEvent(expectedUser)
	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", events.USER_CREATED_TOPIC, mock.MatchedBy(
		func(actualUserEvent events.UserCreatedEvent) bool {
			expectedUserEvent.Id = actualUserEvent.Id
			expectedUserEvent.Password = actualUserEvent.Password
			expectedUserEvent.CreatedAt = actualUserEvent.CreatedAt
			expectedUserEvent.UpdatedAt = actualUserEvent.UpdatedAt
			return assert.Equal(t, expectedUserEvent, actualUserEvent)
		},
	)).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)
	userId, err := userservice.CreateUser(context.Background(), userCreation)

	mockUserStorage.AssertNumberOfCalls(t, "Insert", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 1)

	require.Nil(t, err)
	assert.NotNil(t, userId)
}

// Test when DB fails for some reason while inserting
func TestCreateUserExpectEmptyUserIdWhenInsertFails(t *testing.T) {
	userCreation := testDataUserCreation()

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("Insert", mock.Anything).Return(errors.New("test for failure"))

	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)
	userId, err := userservice.CreateUser(context.Background(), userCreation)

	mockUserStorage.AssertNumberOfCalls(t, "Insert", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 0) // Expect that it will NOT be called

	require.NotNil(t, err)
	assert.Empty(t, userId)
}

// Tests for Retrieving all users
func TestGetUsers(t *testing.T) {
	existingUsers := testDataUsers()
	queryParams := models.GetUserQueryParams{}

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("GetAll", queryParams).Return(existingUsers, nil)

	mockMessageBroker := new(mocks.MockMessageBroker)
	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	actualUsers, err := userservice.GetUsers(context.Background(), queryParams)

	mockUserStorage.AssertNumberOfCalls(t, "GetAll", 1)
	require.Nil(t, err)
	assert.Equal(t, existingUsers, actualUsers)
}

// Test when DB fails for some reason while retrieving
func TestGetUsersExpectNilWhenFetchingFails(t *testing.T) {
	queryParams := models.GetUserQueryParams{}

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("GetAll", queryParams).Return(nil, errors.New("test for failure"))

	mockMessageBroker := new(mocks.MockMessageBroker)
	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	actualUsers, err := userservice.GetUsers(context.Background(), queryParams)

	mockUserStorage.AssertNumberOfCalls(t, "GetAll", 1)
	require.NotNil(t, err)
	assert.Nil(t, actualUsers)
}

// Test for retrieving single user by ID
func TestGetUserById(t *testing.T) {
	existingUser := testDataUser()
	userId := "bf06fc2c-823c-4827-ba08-d3c5fe623708"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("GetById", userId).Return(&existingUser, nil)

	mockMessageBroker := new(mocks.MockMessageBroker)
	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	actualUser, err := userservice.GetUserById(context.Background(), userId)

	mockUserStorage.AssertNumberOfCalls(t, "GetById", 1)
	require.Nil(t, err)
	assert.Equal(t, existingUser, *actualUser)
}

// Test when the user with given ID does not exist in DB
func TestGetUserByIdWhenNonExistingUserId(t *testing.T) {
	userId := "non-existing-user-id"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("GetById", userId).Return(nil, nil)

	mockMessageBroker := new(mocks.MockMessageBroker)
	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	actualUser, err := userservice.GetUserById(context.Background(), userId)

	mockUserStorage.AssertNumberOfCalls(t, "GetById", 1)
	require.Nil(t, err)
	assert.Nil(t, actualUser)
}

// Test when DB fails for some reason while retrieving
func TestGetUserByIdWhenFetchingFails(t *testing.T) {
	userId := "non-existing-user-id"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("GetById", userId).Return(nil, errors.New("test for failure"))

	mockMessageBroker := new(mocks.MockMessageBroker)
	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	actualUser, err := userservice.GetUserById(context.Background(), userId)

	mockUserStorage.AssertNumberOfCalls(t, "GetById", 1)
	require.NotNil(t, err)
	assert.Nil(t, actualUser)
}

// Test for updating a user
func TestUpdateUser(t *testing.T) {
	userToUpdate := testDataUserCreation()
	userId := "bf06fc2c-823c-4827-ba08-d3c5fe623708"

	expectedUser := models.User{
		Id:        userId,
		FirstName: userToUpdate.FirstName,
		LastName:  userToUpdate.LastName,
		Nickname:  userToUpdate.Nickname,
		Email:     userToUpdate.Email,
		Country:   userToUpdate.Country,
	}

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("Update", mock.MatchedBy(
		func(actualUser models.User) bool {
			expectedUser.Password = actualUser.Password
			expectedUser.UpdatedAt = actualUser.UpdatedAt
			return assert.Equal(t, expectedUser, actualUser)
		},
	)).Return(nil)

	expectedUserEvent := events.UserUpdatedEvent(expectedUser)
	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", events.USER_UPDATED_TOPIC, mock.MatchedBy(
		func(actualUserEvent events.UserUpdatedEvent) bool {
			expectedUserEvent.Id = actualUserEvent.Id
			expectedUserEvent.Password = actualUserEvent.Password
			expectedUserEvent.UpdatedAt = actualUserEvent.UpdatedAt
			return assert.Equal(t, expectedUserEvent, actualUserEvent)
		},
	)).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)
	err := userservice.UpdateUser(context.Background(), userId, userToUpdate)

	mockUserStorage.AssertNumberOfCalls(t, "Update", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 1)

	require.Nil(t, err)
}

// Test when DB fails for some reason while updating
func TestUpdateUserWhenUpdateFails(t *testing.T) {
	userToUpdate := testDataUserCreation()
	userId := "doesnotmatter"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("Update", mock.Anything).Return(errors.New("test for failure"))

	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)
	err := userservice.UpdateUser(context.Background(), userId, userToUpdate)

	mockUserStorage.AssertNumberOfCalls(t, "Update", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 0) // Expect that it will NOT be called

	require.NotNil(t, err)
}

// Test for deleting a user by ID
func TestDeleteUser(t *testing.T) {
	userId := "bf06fc2c-823c-4827-ba08-d3c5fe623708"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("DeleteById", userId).Return(nil)

	expectedUserDeletedEvent := events.UserDeletedEvent{Id: userId}
	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", events.USER_DELETED_TOPIC, expectedUserDeletedEvent).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	err := userservice.DeleteUserById(context.Background(), userId)

	mockUserStorage.AssertNumberOfCalls(t, "DeleteById", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 1)
	require.Nil(t, err)
}

// Test when DB fails for some reason while deleting
func TestDeleteUserWhenDeleteFails(t *testing.T) {
	userId := "bf06fc2c-823c-4827-ba08-d3c5fe623708"

	mockUserStorage := new(mocks.MockUserStorage)
	mockUserStorage.On("DeleteById", userId).Return(errors.New("test for db failure"))

	mockMessageBroker := new(mocks.MockMessageBroker)
	mockMessageBroker.On("Publish", events.USER_DELETED_TOPIC, mock.Anything).Return(nil)

	userservice := services.NewUserService(mockUserStorage, mockMessageBroker)

	err := userservice.DeleteUserById(context.Background(), userId)

	mockUserStorage.AssertNumberOfCalls(t, "DeleteById", 1)
	mockMessageBroker.AssertNumberOfCalls(t, "Publish", 0)
	require.NotNil(t, err)
}

// Test Data
func testDataUserCreation() models.UserCreation {
	return models.UserCreation{
		FirstName: "Alice",
		LastName:  "Mustermann",
		Nickname:  "alice",
		Password:  "s3cr3t",
		Email:     "alice@example.com",
		Country:   "UK",
	}
}

func testDataUser() models.User {
	return models.User{
		Id:        "bf06fc2c-823c-4827-ba08-d3c5fe623708",
		FirstName: "Alice",
		LastName:  "Mustermann",
		Nickname:  "alice",
		Password:  "s3cr3t",
		Email:     "alice@example.com",
		Country:   "UK",
	}
}

func testDataUsers() []models.User {
	return []models.User{{
		Id:        "bf06fc2c-823c-4827-ba08-d3c5fe623708",
		Password:  "s3cr3t",
		FirstName: "Alice",
		LastName:  "Mustermann",
		Nickname:  "alice",
		Email:     "alice@example.com",
		Country:   "UK",
		CreatedAt: "2022-06-24T11:58:50Z",
		UpdatedAt: "2022-06-24T11:58:50Z",
	},
		{
			Id:        "fb5cd2c0-67f0-43d1-9c83-a5db0d05c5e8",
			Password:  "s3r10usS4m",
			FirstName: "Samuel",
			LastName:  "Wilson",
			Nickname:  "serious_sam",
			Email:     "sam@example.com",
			Country:   "UK",
			CreatedAt: "2022-06-24T12:12:50Z",
			UpdatedAt: "2022-06-24T12:31:50Z",
		}}
}
