package statemachine

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StateMachineSuite struct {
	suite.Suite
}

func TestStateMachineSuite(t *testing.T) {
	suite.Run(t, new(StateMachineSuite))
}

type testInputEvent struct {
	kind InputEventKind
}

func (e testInputEvent) Kind() InputEventKind {
	return e.kind
}

func (s *StateMachineSuite) TestNewStateMachine() {
	var notifications []StateTransition[string]
	recorderNotifier := func(transition StateTransition[string]) {
		notifications = append(notifications, transition)
	}

	states := map[StateId]*State[string]{
		"state1": {
			Id: "state1",
			TransitionHandler: func(event InputEvent, data string) (StateId, InternalDataMutator[string], error) {
				switch event.Kind() {
				case "event-a":
					// On event-a, move to state 2 and don't mutate data
					return "state2", nil, nil
				case "event-b":
					// On event-b, move to state 3 and change the internal data to "bongo!"
					return "state3", func(s *string) {
						*s = "bongo!"
					}, nil
				default:
					// event-c is not allowed
					return "", nil, errors.New("invalid event")
				}
			},
		},
		"state2": {
			Id: "state2",
			TransitionHandler: func(event InputEvent, data string) (StateId, InternalDataMutator[string], error) {
				switch event.Kind() {
				case "event-a":
					// On event-a, stay in state 2 and add an exclamation mark to the data
					return "state2", func(s *string) {
						*s += "!"
					}, nil
				case "event-b":
					// On event-b, move to state 3 without changing the data
					return "state3", nil, nil
				case "event-c":
					// On event-c, move to state 1 without changing the data
					return "state1", nil, nil
				default:
					return "", nil, errors.New("invalid event")
				}
			},
		},
		"state3": {
			Id: "state3",
			TransitionHandler: func(event InputEvent, data string) (StateId, InternalDataMutator[string], error) {
				switch event.Kind() {
				case "event-b":
					// On event-b, move to state 2 without changing the data
					return "state2", nil, nil
				case "event-c":
					// On event-c, move to state1 and add " LOOP" to the data
					return "state1", func(s *string) {
						*s += " LOOP"
					}, nil
				default:
					// event-a is not allowed
					return "", nil, errors.New("invalid event")
				}
			},
		},
	}

	sm := NewStateMachine(states, recorderNotifier, "state1", "a")
	s.NotNil(sm)

	// Initially, the state machine is in state1 with data "a"
	s.Equal(StateId("state1"), sm.CurrentState)
	s.Equal("a", sm.InternalData)

	// event-c is invalid in state1
	err := sm.HandleEvent(testInputEvent{"event-c"})
	s.Error(err)

	// event-a should move us to state2, with no change to the data
	err = sm.HandleEvent(testInputEvent{"event-a"})
	s.NoError(err)
	s.Equal(StateId("state2"), sm.CurrentState)
	s.Equal("a", sm.InternalData)

	// event-a should keep us in state2, with an exclamation mark added to the data
	err = sm.HandleEvent(testInputEvent{"event-a"})
	s.NoError(err)
	s.Equal(StateId("state2"), sm.CurrentState)
	s.Equal("a!", sm.InternalData)

	// Repeating event-a should keep us in state2, with another exclamation mark added to the data
	err = sm.HandleEvent(testInputEvent{"event-a"})
	s.NoError(err)
	s.Equal(StateId("state2"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-c should move us to state1, with no data change
	err = sm.HandleEvent(testInputEvent{"event-c"})
	s.NoError(err)
	s.Equal(StateId("state1"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-a should move us back to state2, with no data change
	err = sm.HandleEvent(testInputEvent{"event-a"})
	s.NoError(err)
	s.Equal(StateId("state2"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-b should move us to state3, with no data change
	err = sm.HandleEvent(testInputEvent{"event-b"})
	s.NoError(err)
	s.Equal(StateId("state3"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-a is not valid in state3
	err = sm.HandleEvent(testInputEvent{"event-a"})
	s.Error(err)

	// event-b should move us back to state2, with no data change
	err = sm.HandleEvent(testInputEvent{"event-b"})
	s.NoError(err)
	s.Equal(StateId("state2"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-b should move us back to state3, with no data change
	err = sm.HandleEvent(testInputEvent{"event-b"})
	s.NoError(err)
	s.Equal(StateId("state3"), sm.CurrentState)
	s.Equal("a!!", sm.InternalData)

	// event-c should move us back to state1, with " LOOP" added to the data
	err = sm.HandleEvent(testInputEvent{"event-c"})
	s.NoError(err)
	s.Equal(StateId("state1"), sm.CurrentState)
	s.Equal("a!! LOOP", sm.InternalData)

	// event-b should move us to state3, and replace our data with "bongo!"
	err = sm.HandleEvent(testInputEvent{"event-b"})
	s.NoError(err)
	s.Equal(StateId("state3"), sm.CurrentState)
	s.Equal("bongo!", sm.InternalData)

	// Validate transition sequence
	expectedSequence := []StateTransition[string]{
		{"state1", "state2"},
		{"state2", "state2"},
		{"state2", "state2"},
		{"state2", "state1"},
		{"state1", "state2"},
		{"state2", "state3"},
		{"state3", "state2"},
		{"state2", "state3"},
		{"state3", "state1"},
		{"state1", "state3"},
	}
	s.Equal(expectedSequence, notifications)
}
