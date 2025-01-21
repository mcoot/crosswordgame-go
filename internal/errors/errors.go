package errors

import (
	"errors"
	"fmt"
)

type GameErrorKind string

const (
	GameErrorInvalidInput        GameErrorKind = "invalid_input"
	GameErrorNotFound            GameErrorKind = "not_found"
	GameErrorInvalidAction       GameErrorKind = "invalid_action"
	GameErrorUnexpectedGameLogic GameErrorKind = "unexpected_game_logic_error"
)

type GameError interface {
	error
	HTTPCode() int
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

func IsNotFoundError(err error) bool {
	ge, ok := AsGameError(err)
	return ok && ge.Kind() == GameErrorNotFound
}

type InvalidInputError struct {
	ErrMessage string
}

func (e *InvalidInputError) HTTPCode() int {
	return 400
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

func (e *NotFoundError) HTTPCode() int {
	return 404
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

type InvalidActionError struct {
	Action string
	Reason string
}

func (e *InvalidActionError) HTTPCode() int {
	return 400
}

func (e *InvalidActionError) Kind() GameErrorKind {
	return GameErrorInvalidAction
}

func (e *InvalidActionError) Message() string {
	return fmt.Sprintf("invalid action \"%s\": %s", e.Action, e.Reason)
}

func (e *InvalidActionError) Error() string {
	return e.Message()
}

type UnexpectedGameLogicError struct {
	ErrMessage string
}

func (e *UnexpectedGameLogicError) HTTPCode() int {
	return 500
}

func (e *UnexpectedGameLogicError) Kind() GameErrorKind {
	return GameErrorUnexpectedGameLogic
}

func (e *UnexpectedGameLogicError) Message() string {
	return e.ErrMessage
}

func (e *UnexpectedGameLogicError) Error() string {
	return e.Message()
}
