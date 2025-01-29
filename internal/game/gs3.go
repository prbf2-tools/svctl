package game

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"regexp"
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
	port, err := s.gs3Port()
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("udp", net.JoinHostPort("localhost", port))
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

func (s *Server) gs3Port() (string, error) {
	content, err := os.ReadFile(filepath.Join(s.Path, "mods/pr", "settings/serversettings.con"))
	if err != nil {
		return "", err
	}

	regex := `sv.gameSpyPort\s*(\d+)`
	return extractWithPattern(string(content), regex), nil
}

func extractWithPattern(content string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
