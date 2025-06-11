package http

import "net/http"

type ErrorWrapper struct {
	StatusCode int
	Message    string
}

func (e *ErrorWrapper) Error() string {
	return e.Message
}

func (e *ErrorWrapper) HTTPStatus() int {
	return e.StatusCode
}

func NewInternalServerError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusInternalServerError,
		Message:    errorMessage,
	}
}

func NewBadRequestError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusBadRequest,
		Message:    errorMessage,
	}
}

func NewNotFoundError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusNotFound,
		Message:    errorMessage,
	}
}

func NewUnauthorizedError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusUnauthorized,
		Message:    errorMessage,
	}
}

func NewConflictError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusConflict,
		Message:    errorMessage,
	}
}

func NewUnprocessableEntityError(errorMessage string) error {
	return &ErrorWrapper{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    errorMessage,
	}
}
