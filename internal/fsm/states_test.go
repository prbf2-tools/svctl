package fsm

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type StatesSuite struct {
	suite.Suite
}

func (s *StatesSuite) TestStateErrored() {
	err := errors.New("test error")
	state := NewStateErrored(err)

	nextState, stateErr := state.EventHandler(EventReset, nil)
	s.IsType(&StateStopped{}, nextState)
	s.Equal(err, stateErr)
}

func (s *StatesSuite) TestStateStopped() {
	ctrl := gomock.NewController(s.T())
	gameServerMock := NewMockGameServer(ctrl)

	state := NewStateStopped()
	fsm := New(gameServerMock, state)

	gameServerMock.EXPECT().Start().Return(nil)

	nextState, stateErr := state.EventHandler(EventStart, fsm)
	s.IsType(&StateRunning{}, nextState)
	s.Nil(stateErr)

	nextState, stateErr = state.EventHandler(EventStop, fsm)
	s.Nil(nextState)
	s.Equal(ErrEventNotAllowed, stateErr)
}

func (s *StatesSuite) TestStateRunning() {
	ctrl := gomock.NewController(s.T())
	gameServerMock := NewMockGameServer(ctrl)

	state := NewStateRunning(nil)
	state.renderInterval = time.Second

	fsm := New(gameServerMock, state)

	s.Run("Render should be called", func() {
		gameServerMock.EXPECT().IsRunning().Return(true).AnyTimes()
		renderCalled := false
		gameServerMock.EXPECT().Render().Return(nil).Do(func() {
			renderCalled = true
		})

		state.OnEnter(fsm)

		s.Eventually(func() bool {
			return renderCalled
		}, 3*time.Second, 100*time.Millisecond, "Render should be called")

		state.OnExit()
	})

	s.Run("Stop should be called", func() {
		gameServerMock.EXPECT().Stop().Return(nil)

		nextState, stateErr := state.EventHandler(EventStop, fsm)
		s.IsType(&StateStopped{}, nextState)
		s.Nil(stateErr)
	})

	s.Run("Succesfull restart", func() {
		gameServerMock.EXPECT().Stop().Return(nil)

		nextState, stateErr := state.EventHandler(EventRestart, fsm)
		s.IsType(&StateRestarting{}, nextState)
		s.Nil(stateErr)
	})

	s.Run("Failed restart", func() {
		err := errors.New("test error")
		gameServerMock.EXPECT().Stop().Return(err)

		nextState, stateErr := state.EventHandler(EventRestart, fsm)
		s.IsType(&StateErrored{}, nextState)
		s.Equal(err, stateErr)
	})

	s.Run("Not allowed event", func() {
		nextState, stateErr := state.EventHandler(EventReset, fsm)
		s.Nil(nextState)
		s.Equal(ErrEventNotAllowed, stateErr)
	})
}

func (s *StatesSuite) TestStateRestarting() {
	ctrl := gomock.NewController(s.T())
	gameServerMock := NewMockGameServer(ctrl)

	state := NewStateRestarting(NewRestartCounter(3))
	fsm := New(gameServerMock, state)

	s.Run("Succesfull restart", func() {
		gameServerMock.EXPECT().Start().Return(nil)
		state.OnEnter(fsm)

		s.IsType(&StateRunning{}, fsm.desiredState)
	})

	s.Run("Failed restart", func() {
		err := errors.New("test error")
		gameServerMock.EXPECT().Start().Return(err)
		state.OnEnter(fsm)

		s.IsType(&StateErrored{}, fsm.desiredState)
		s.Equal(err, fsm.desiredState.(*StateErrored).Err)
	})

	s.Run("Max restarts reached", func() {
		state.OnEnter(fsm)
		s.IsType(&StateErrored{}, fsm.desiredState)
		s.Equal(ErrMaxRestartsReached, fsm.desiredState.(*StateErrored).Err)
	})
}

func TestStatesSuite(t *testing.T) {
	suite.Run(t, new(StatesSuite))
}
