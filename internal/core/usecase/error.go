package usecase

import (
	"errors"
	"fmt"
)

type UserError interface {
	error
	UserMessage() string
}

type Error struct {
	message string
	err     error
}

func NewError(message string, err error) *Error {
	return &Error{
		message: message,
		err:     err,
	}
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) UserMessage() string {
	return e.message
}

func MessageFor(err error) string {
	if err == nil {
		return ""
	}

	var userErr UserError
	if ok := errors.As(err, &userErr); ok {
		return userErr.UserMessage()
	}

	return "An unexpected error occurred."
}
