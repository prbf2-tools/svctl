package fsm

type StateStopped struct {
	baseState
}

func NewStateStopped() *StateStopped {
	return &StateStopped{}
}

func (s *StateStopped) EventHandler(event Event, fsm *FSM) (State, error) {
	const op = "StateStopped.EventHandler"
	log := fsm.Log.With("op", op)
	log.Debug("Received event", "event", event)

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
