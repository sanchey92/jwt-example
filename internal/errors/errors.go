package appError

import (
	"errors"
	"net/http"
)

var (
	ErrMissingEnvVars = errors.New("missing required environment variables")
	ErrInvalidInput   = errors.New("invalid input data")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
)

type ApiError struct {
	StatusCode int
	Message    string
}

func NewApiError(statusCode int, err error) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

func (e *ApiError) Error() string {
	return e.Message
}

func BadRequest(err error) *ApiError {
	return NewApiError(http.StatusBadRequest, err)
}

func Unauthorized(err error) *ApiError {
	return NewApiError(http.StatusUnauthorized, err)
}

func Forbidden(err error) *ApiError {
	return NewApiError(http.StatusForbidden, err)
}

func InternalServer(err error) *ApiError {
	return NewApiError(http.StatusInternalServerError, err)
}
