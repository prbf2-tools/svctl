package fsm

import (
	_ "go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -source=fsm.go -destination=mock_game_server.go -package=fsm GameServer
