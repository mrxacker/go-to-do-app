package errors

import "errors"

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidIdentifier = errors.New("invalid identifier")
)
