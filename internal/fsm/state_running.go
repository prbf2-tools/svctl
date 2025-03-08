package fsm

import (
	"context"
	"time"
)

const (
	maxRestarts    = 5
	renderInterval = time.Minute
)

type StateRunning struct {
	counter        *restartCounter
	renderInterval time.Duration

	cancel context.CancelFunc
}

func NewStateRunning(counter *restartCounter) *StateRunning {
	if counter == nil {
		counter = NewRestartCounter(maxRestarts)
	}

	return &StateRunning{
		counter:        counter,
		renderInterval: renderInterval,
	}
}

func (s *StateRunning) OnEnter(fsm *FSM) {
	const op = "StateRunning.OnEnter"
	log := fsm.Log.With("op", op)

	log.Debug("Server running")

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	ticker := time.NewTicker(s.renderInterval)

	sv := fsm.Server()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("Running loop cancelled")
				ticker.Stop()
				return
			case <-ticker.C:
				log.Debug("Rendering templates (reloadable)")
				err := sv.Render(true)
				if err != nil {
					log.Error("Failed to render templates", "err", err)
				}
			default:
				isRunning, err := sv.IsRunning()
				if err != nil {
					log.Error("Failed to check if server is running", "err", err)
					continue
				}

				if !isRunning {
					log.Error("Server not running, attempting restart")
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
	const op = "StateStopped.EventHandler"
	log := fsm.Log.With("op", op)
	log.Debug("Received event", "event", event)

	switch event {
	case EventStop:
		s.cancel()
		if err := fsm.Server().Stop(); err != nil {
			log.Error("Failed to stop server", "err", err)
			return NewStateErrored(err), err
		}

		return NewStateStopped(), nil
	case EventRestart:
		if err := fsm.Server().Stop(); err != nil {
			log.Error("Failed to restart server", "err", err)
			return NewStateErrored(err), err
		}

		return NewStateRestarting(s.counter), nil
	default:
		return nil, ErrEventNotAllowed
	}
}
