package fsm

import (
	"context"
	"log/slog"
	"time"
)

type GameServer interface {
	Start() error
	Stop() error
	IsRunning() bool
	Render() error
}

type FSM struct {
	server GameServer
	log    *slog.Logger

	currentState State
	desiredState State

	cancel context.CancelFunc
}

func New(server GameServer, log *slog.Logger, initialState State) *FSM {
	fsm := FSM{
		currentState: &baseState{},
		desiredState: initialState,
		server:       server,
		log:          log,
	}

	fsm.Transition()
	go fsm.Run()

	return &fsm
}

func (f *FSM) Server() GameServer {
	return f.server
}

func (f *FSM) ChangeState(state State) {
	f.desiredState = state
}

func (f *FSM) Event(event Event) error {
	println("event", event)
	if f.currentState == nil {
		return nil
	}

	// Error is for the user, state is for the FSM
	nextState, err := f.currentState.EventHandler(event, f)
	if nextState != nil {
		f.desiredState = nextState
	}

	return err
}

func (f *FSM) Run() {
	if f.cancel != nil {
		f.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel

	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.Transition()
		}
	}
}

func (f *FSM) Transition() {
	if f.desiredState != f.currentState {
		if f.currentState != nil {
			f.currentState.OnExit()
		}

		f.currentState = f.desiredState
		f.currentState.OnEnter(f)
	}
}
