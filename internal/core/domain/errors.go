package domain

import (
	"errors"
	"net/http"
)

type Error struct {
	status  int
	code    string
	message string
	err     error
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) Code() string {
	return e.code
}

func NewError(status int, code, message string, err error) *Error {
	return &Error{
		status:  status,
		code:    code,
		message: message,
		err:     err,
	}
}

func ValidationError(message string, err error) *Error {
	return NewError(http.StatusBadRequest, "invalid_request", message, err)
}

func NotFoundError(resource string, err error) *Error {
	return NewError(http.StatusNotFound, "not_found", resource+" not found", err)
}

func ConflictError(message string, err error) *Error {
	return NewError(http.StatusConflict, "conflict", message, err)
}

func ServiceUnavailableError(message string, err error) *Error {
	return NewError(http.StatusServiceUnavailable, "service_unavailable", message, err)
}

func InternalError(message string, err error) *Error {
	return NewError(http.StatusInternalServerError, "internal_error", message, err)
}

func StatusCode(err error) int {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Status()
	}

	return http.StatusInternalServerError
}

func ErrorCode(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code()
	}

	return "internal_error"
}

func ErrorMessage(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Error()
	}

	return "internal server error"
}
