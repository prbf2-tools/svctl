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
	if s.counter != nil {
		s.counter.Increment()
		if s.counter.LimitReached() {
			fsm.ChangeState(NewStateErrored(ErrMaxRestartsReached))
			return
		}
	}

	err := fsm.Server().Start()
	if err != nil {
		fsm.ChangeState(NewStateErrored(err))
		return
	}

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
