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

func (s *StateErrored) OnEnter(fsm *FSM) {
	const op = "StateErrored.OnEnter"
	log := fsm.Log.With("op", op)

	log.Error("Server errored", "err", s.Err)
}

func (s *StateErrored) EventHandler(event Event, fsm *FSM) (State, error) {
	const op = "StateErrored.EventHandler"
	log := fsm.Log.With("op", op)
	log.Debug("Received event", "event", event)

	switch event {
	case EventReset:
		return NewStateStopped(), s.Err
	default:
		return nil, ErrEventNotAllowed
	}
}
