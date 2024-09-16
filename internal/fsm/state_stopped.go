package fsm

type StateStopped struct {
	baseState
}

func NewStateStopped() *StateStopped {
	return &StateStopped{}
}

func (s *StateStopped) EventHandler(event Event, fsm *FSM) (State, error) {
	switch event {
	case EventStart:
		err := fsm.Server().Start()
		if err != nil {
			return NewStateErrored(err), err
		}
		return NewStateRunning(nil), nil
	default:
		return nil, ErrEventNotAllowed
	}
}
