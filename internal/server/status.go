package server

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
	// config, err := s.Config()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// host := "127.0.0.1"
	// if config.Game.ExternalIP != "" {
	// 	host = config.Game.ServerIP
	// }
	//
	// serverAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, strconv.Itoa(config.Game.GamespyPort)))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// udp, err := net.DialUDP("udp", nil, serverAddr)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// defer udp.Close()
	//
	// c := gamespy3.Status(context.Background(), udp)
	//
	// resp, err := c.ServerInfoB(context.Background())
	// if err != nil {
	// 	return nil, err
	// }
	//
	// status := &Status{
	// 	Hostname:   resp.Header.Hostname,
	// 	Port:       strconv.Itoa(resp.Header.HostPort),
	// 	MapName:    resp.Header.Map.Name,
	// 	GameMode:   resp.Header.Map.GameMode,
	// 	MapSize:    resp.Header.Map.Layer,
	// 	NumPlayers: resp.Header.NumPlayers,
	// 	MaxPlayers: resp.Header.MaxPlayers,
	// }
	//
	// return status, nil
	// TODO: Implement actual status retrieval
	return &Status{}, nil
}
