package client

import "net"

// InitTCP struct contains data for init player
type InitTCP struct {
	Method     string     `json:"method"`
	Nickname   string     `json:"nickname"`
	Resolution Resolution `json:"resolution"`
}

// InitUDP struct contains data for init player from UDP
type InitUDP struct {
	Method   string `json:"method"`
	Nickname string `json:"nickname"`
}

// Disconnect struct contains data for disconnect player
type Disconnect struct {
	Method   string `json:"method"`
	Nickname string `json:"nickname"`
}

// Resolution struct contains resolution
type Resolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Position - position of cursor
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Mouse - position of cursor
type Mouse struct {
	Nickname  string   `json:"nickname"`
	Method    string   `json:"method"`
	IsClicked bool     `json:"is_clicked"`
	Position  Position `json:"position"`
}

// Keyboard - events from player's keyboard
type Keyboard struct {
	Nickname string   `json:"nickname"`
	Method   string   `json:"method"`
	Keys     []string `json:"keys"`
}

// Response - response struct's data
type Response struct {
	Method  string `json:"method"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// UDPNetData struct contains UDP's connection's data
type UDPNetData struct {
	Addr     *net.UDPAddr
	Nickname string
}

// TCPNetData struct contains TCP's connection's data
type TCPNetData struct {
	Addr     net.Addr
	Nickname string
}
