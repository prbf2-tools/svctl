package server

import (
	"context"
	"net"
	"strconv"

	"github.com/emilekm/go-prbf2/gamespy3"
)

type Status struct {
	Hostname   string `json:"hostname"`
	Port       string `json:"hostport"`
	GameMode   string `json:"gamemode"`
	MapName    string `json:"mapname"`
	MapSize    int    `json:"mapsize"`
	NumPlayers int    `json:"numplayers"`
	MaxPlayers int    `json:"maxplayers"`
}

func (s *Server) Status() (*Status, error) {
	config, err := s.Config()
	if err != nil {
		return nil, err
	}

	host := "127.0.0.1"
	if config.Game.ExternalIP != "" {
		host = config.Game.ServerIP
	}

	conn, err := net.Dial("udp", net.JoinHostPort(host, config.Game.GamespyPort))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	resp, err := gamespy3.Status(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	status := &Status{
		Hostname:   resp.Hostname,
		Port:       strconv.Itoa(resp.HostPort),
		MapName:    resp.Map.Name,
		GameMode:   resp.Map.GameMode,
		MapSize:    resp.Map.Layer,
		NumPlayers: resp.NumPlayers,
		MaxPlayers: resp.MaxPlayers,
	}

	return status, nil
}
