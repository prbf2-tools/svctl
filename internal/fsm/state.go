package fsm

import "errors"

// Stopped
// Adopting
// Starting
// Running
// Stopping
// Restarting
// AutoRestarting
// Exited

var (
	ErrEventNotAllowed = errors.New("event not allowed in current state")
)

type State interface {
	OnEnter(*FSM)
	OnExit()
	EventHandler(Event, *FSM) (State, error)
}

type baseState struct{}

func (s *baseState) OnEnter(_ *FSM) {}
func (s *baseState) OnExit()        {}
func (s *baseState) EventHandler(_ Event, _ *FSM) (State, error) {
	return nil, nil
}
