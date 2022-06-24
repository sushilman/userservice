package apiuser_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/sushilman/userservice/api/apiuser"
	"github.com/sushilman/userservice/mocks"
)

func TestDeleteUserByIdHandler(t *testing.T) {
	w := httptest.NewRecorder()
	gin.CreateTestContext(w)

	mockUserService := new(mocks.MockUserService)

	mockUserService.On("DeleteUserById", "1").Return(errors.New("random error"))

	apiuser.DeleteUserByIdHandler(mockUserService)

	assert.Equal(t, 204, w.Code)
}
