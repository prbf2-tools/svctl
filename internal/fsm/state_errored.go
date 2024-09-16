package fsm

type StateErrored struct {
	baseState
	Err error
}

func NewStateErrored(err error) *StateErrored {
	return &StateErrored{
		Err: err,
	}
}

func (s *StateErrored) EventHandler(event Event, fsm *FSM) (State, error) {
	switch event {
	case EventReset:
		return NewStateStopped(), s.Err
	default:
		return nil, ErrEventNotAllowed
	}
}
