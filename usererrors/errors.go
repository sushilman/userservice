package usererrors

import (
	"net/http"

	"github.com/sushilman/userservice/models"
)

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "resource not found"
}

// API error responses in JSON format

func NewNotFoundErrorResponse(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func NewBadRequestErrorResponse(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

func NewInternalServerError(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}
