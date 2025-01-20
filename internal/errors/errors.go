package errors

import (
	"errors"
	"fmt"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
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

type InvalidActionError struct {
	PlayerId playertypes.PlayerId
	Action   string
	Reason   string
}

func (e *InvalidActionError) Kind() GameErrorKind {
	return GameErrorInvalidAction
}

func (e *InvalidActionError) Message() string {
	return fmt.Sprintf("invalid action \"%s\" for player %s: %s", e.Action, e.PlayerId, e.Reason)
}

func (e *InvalidActionError) Error() string {
	return e.Message()
}

type UnexpectedGameLogicError struct {
	ErrMessage string
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
