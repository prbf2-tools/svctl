package daemon

import (
	"errors"

	"github.com/prbf2-tools/svctl/internal/game"
	"github.com/prbf2-tools/svctl/internal/server"
)

type Event string

const (
	EventStart   Event = "start"
	EventStop    Event = "stop"
	EventRestart Event = "restart"
	EventReset   Event = "reset"
)

type State int

const (
	StateStopped State = iota
	StateStarting
	StateRunning
	StateRestarting
	StateStopping
	StateExited
)

type eventHandler func(*server.Server) (State, error)

type transition struct {
	transitionState State
	handler         eventHandler
}

var transitionTable = map[State]map[Event]transition{
	StateStopped: {
		EventStart: {
			transitionState: StateStarting,
			handler:         startServer,
		},
	},
	StateRunning: {
		EventStop: {
			transitionState: StateStopping,
			handler:         stopServer,
		},
		EventRestart: {
			transitionState: StateRestarting,
			handler:         restartServer,
		},
	},
	StateExited: {
		EventReset: {
			transitionState: StateStopped,
			handler:         resetServer,
		},
	},
}

func startServer(s *server.Server) (State, error) {
	err := s.ApplyPatches()
	if err != nil {
		return StateExited, err
	}

	err = s.Render(false)
	if err != nil {
		return StateExited, err
	}

	err = s.Start()
	if err != nil && errors.Is(err, game.ErrAlreadyRunning) {
		return StateExited, err
	}

	return StateRunning, nil
}

func stopServer(s *server.Server) (State, error) {
	err := s.Stop()
	if err != nil && !errors.Is(err, game.ErrNotRunning) {
		return StateRunning, err
	}

	return StateStopped, nil
}

func restartServer(s *server.Server) (State, error) {
	newState, err := stopServer(s)
	if err != nil {
		return newState, err
	}

	newState, err = startServer(s)
	if err != nil {
		return newState, err
	}

	return newState, nil
}

func resetServer(s *server.Server) (State, error) {
	return StateStopped, nil
}
