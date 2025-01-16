package types

import (
	"errors"
	"fmt"
)

type GameErrorKind string

const (
	GameErrorInvalidInput GameErrorKind = "invalid_input"
	GameErrorNotFound     GameErrorKind = "not_found"
)

type GameError interface {
	error
	Kind() GameErrorKind
	Message() string
}

func AsGameError(err error) (GameError, bool) {
	if err == nil {
		return nil, false
	}
	var ge GameError
	if errors.As(err, &ge) {
		return ge, true
	}
	return nil, false
}

type InvalidInputError struct {
	ErrMessage string
}

func (e *InvalidInputError) Kind() GameErrorKind {
	return GameErrorInvalidInput
}

func (e *InvalidInputError) Message() string {
	return e.ErrMessage
}

func (e *InvalidInputError) Error() string {
	return e.Message()
}

type NotFoundError struct {
	ObjectKind string
	ObjectID   interface{}
}

func (e *NotFoundError) Kind() GameErrorKind {
	return GameErrorNotFound
}

func (e *NotFoundError) Message() string {
	return fmt.Sprintf("%s \"%v\" not found", e.ObjectKind, e.ObjectID)
}

func (e *NotFoundError) Error() string {
	return e.Message()
}
