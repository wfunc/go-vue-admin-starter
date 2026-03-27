package errorx

import (
	"errors"
	"net/http"
)

type Error struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func Wrap(status int, code, message string, err error) *Error {
	return &Error{Status: status, Code: code, Message: message, Err: err}
}

func BadRequest(message string) *Error { return New(http.StatusBadRequest, "bad_request", message) }
func Unauthorized(message string) *Error { return New(http.StatusUnauthorized, "unauthorized", message) }
func Forbidden(message string) *Error { return New(http.StatusForbidden, "forbidden", message) }
func NotFound(message string) *Error { return New(http.StatusNotFound, "not_found", message) }
func Conflict(message string) *Error { return New(http.StatusConflict, "conflict", message) }
func Internal(message string) *Error { return New(http.StatusInternalServerError, "internal_error", message) }

func From(err error) *Error {
	if err == nil {
		return nil
	}
	var target *Error
	if errors.As(err, &target) {
		return target
	}
	return Wrap(http.StatusInternalServerError, "internal_error", "internal server error", err)
}
