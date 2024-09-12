package fsm

import (
	"testing"

	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type FSMSuite struct {
	suite.Suite
}

func (s *FSMSuite) TestStartStop() {
	ctrl := gomock.NewController(s.T())
	gameServerMock := NewMockGameServer(ctrl)

	state := NewStateStopped()
	fsm := New(gameServerMock, state)

	gameServerMock.EXPECT().Start().Return(nil)

	err := fsm.Event(EventStart)
	s.Nil(err)
	s.IsType(&StateRunning{}, fsm.desiredState)

	fsm.Transition()

	s.IsType(&StateRunning{}, fsm.currentState)

	gameServerMock.EXPECT().Stop().Return(nil)

	err = fsm.Event(EventStop)
	s.Nil(err)
	s.IsType(&StateStopped{}, fsm.desiredState)

	fsm.Transition()

	s.IsType(&StateStopped{}, fsm.currentState)
}

func TestFSMSuite(t *testing.T) {
	suite.Run(t, new(FSMSuite))
}
