package fsm

import (
	"context"
	"time"
)

const (
	MaxRestarts = 5
)

type StateRunning struct {
	baseState
	counter        *restartCounter
	renderInterval time.Duration

	cancel context.CancelFunc
}

func NewStateRunning(counter *restartCounter) *StateRunning {
	if counter == nil {
		counter = NewRestartCounter(MaxRestarts)
	}

	return &StateRunning{
		counter:        counter,
		renderInterval: time.Minute,
	}
}

func (s *StateRunning) OnEnter(fsm *FSM) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	ticker := time.NewTicker(s.renderInterval)

	sv := fsm.Server()

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				fsm.Server().Render()
			default:
				if !sv.IsRunning() {
					fsm.ChangeState(NewStateRestarting(s.counter))
					ticker.Stop()
					cancel()
					return
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}

func (s *StateRunning) OnExit() {
	s.cancel()
}

func (s *StateRunning) EventHandler(event Event, fsm *FSM) (State, error) {
	switch event {
	case EventStop:
		if err := fsm.Server().Stop(); err != nil {
			return NewStateErrored(err), err
		}

		return NewStateStopped(), nil
	case EventRestart:
		if err := fsm.Server().Stop(); err != nil {
			return NewStateErrored(err), err
		}

		return NewStateRestarting(s.counter), nil
	default:
		return nil, ErrEventNotAllowed
	}
}
