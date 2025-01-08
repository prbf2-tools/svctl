package fsm

import (
	"context"
	"errors"
	"time"
)

var (
	ErrMaxRestartsReached = errors.New("max restarts reached")
)

type StateRestarting struct {
	baseState
	counter *restartCounter
}

func NewStateRestarting(counter *restartCounter) *StateRestarting {
	return &StateRestarting{
		counter: counter,
	}
}

func (s *StateRestarting) OnEnter(fsm *FSM) {
	const op = "StateRestarting.OnEnter"
	log := fsm.Log.With("op", op)

	if s.counter != nil {
		s.counter.Increment()
		if s.counter.LimitReached() {
			log.Error("Max restarts limit reached")
			fsm.ChangeState(NewStateErrored(ErrMaxRestartsReached))
			return
		}
	}

	log.Info("Rendering templates")
	err := fsm.Server().Render(false)
	if err != nil {
		log.Error("Failed to render templates", "err", err)
	}

	log.Info("Restarting server")

	err = fsm.Server().Start()
	if err != nil {
		log.Error("Failed to restart server", "err", err)
		fsm.ChangeState(NewStateErrored(err))
		return
	}

	log.Info("Server successfully restarted")
	fsm.ChangeState(NewStateRunning(s.counter))
}

type restartCounter struct {
	maxRestarts uint8
	numRestarts uint8
	timer       *time.Timer
	cancel      context.CancelFunc
}

func NewRestartCounter(maxRestarts uint8) *restartCounter {
	return &restartCounter{
		maxRestarts: maxRestarts,
	}
}

func (r *restartCounter) restartTimer() {
	if r.timer != nil {
		r.cancel()
		r.timer.Stop()
		r.timer.Reset(time.Minute)
	} else {
		r.timer = time.NewTimer(time.Minute)
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-r.timer.C:
			r.Reset()
		}
	}()
}

func (r *restartCounter) Increment() {
	r.numRestarts++
	r.restartTimer()
}

func (r *restartCounter) Reset() {
	if r.timer != nil {
		r.cancel()
		r.timer.Stop()
	}

	r.numRestarts = 0
}

func (r *restartCounter) LimitReached() bool {
	return r.numRestarts >= r.maxRestarts
}
