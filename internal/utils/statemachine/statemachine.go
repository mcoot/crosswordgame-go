package statemachine

import "fmt"

type InputEventKind string

// InputEvent is an event that a state machine can handle
// It exposes a Kind method for discriminating the event kind, but the specific list of events will depend
// on the state machine
type InputEvent interface {
	Kind() InputEventKind
}

type StateId string

// InternalDataMutator is a function type that mutates the internal data of the state machine
// We return a mutating function rather than mutate inside the transition handler
// because we want the transition handler to be pure, but don't want to copy the internal data constantly
type InternalDataMutator[T any] func(*T)

// TransitionHandler is a function type that represents an edge in the state machine
// Given an event and the current state of the internal data, it should return:
// - the next state
// - a mutator function that will be applied to the internal data during transition
// - an error if the handler cannot apply any transition (even a noop), in which case nothing will be applied
type TransitionHandler[T any] func(event InputEvent, internalData T) (StateId, InternalDataMutator[T], error)

type StateTransition[T any] struct {
	FromState StateId
	ToState   StateId
}

type TransitionNotifier[T any] func(transition StateTransition[T])

type State[T any] struct {
	Id                StateId
	TransitionHandler TransitionHandler[T]
}

type StateMachine[T any] struct {
	States             map[StateId]*State[T]
	TransitionNotifier TransitionNotifier[T]
	CurrentState       StateId
	InternalData       T
}

func NewStateMachine[T any](
	states map[StateId]*State[T],
	transitionNotifier TransitionNotifier[T],
	initialState StateId,
	initialInternalData T,
) *StateMachine[T] {
	return &StateMachine[T]{
		States:             states,
		TransitionNotifier: transitionNotifier,
		CurrentState:       initialState,
		InternalData:       initialInternalData,
	}
}

func (sm *StateMachine[T]) HandleEvent(event InputEvent) error {
	oldState := sm.CurrentState

	state, ok := sm.States[oldState]
	if !ok {
		return fmt.Errorf("state %s not found", oldState)
	}

	nextState, mutator, err := state.TransitionHandler(event, sm.InternalData)
	if err != nil {
		return err
	}

	sm.CurrentState = nextState
	if mutator != nil {
		mutator(&sm.InternalData)
	}

	if sm.TransitionNotifier != nil {
		sm.TransitionNotifier(StateTransition[T]{
			FromState: oldState,
			ToState:   nextState,
		})
	}

	return nil
}
