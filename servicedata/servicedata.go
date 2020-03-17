package servicedata

import (
	"net"

	"github.com/dm1trypon/game-server-golang/models/base"
	"github.com/dm1trypon/game-server-golang/models/block"
	"github.com/dm1trypon/game-server-golang/models/bullet"
	"github.com/dm1trypon/game-server-golang/models/config"
	"github.com/dm1trypon/game-server-golang/models/effect"
	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/dm1trypon/game-server-golang/models/scene"
)

// ConnectedData - contains data on connections to TCP and UDP
// servers and the time to disconnect.
type ConnectedData struct {
	TCPConnect net.Conn
	UDPConnect net.UDPAddr
	TimeDisc   int
	Nickname   string
}

// GameConfig - game service config
var GameConfig config.GameConfig

// ConnectedClients - connected clients to the TCP and UDP servers with time to break
// the connection in case of incorrect identification.
var ConnectedClients []ConnectedData

// UDPClients - connected client's addresses to the UDP server.
var UDPClients map[net.Conn]net.UDPAddr

// Base - the structure of all game objects.
var Base base.Base

// PlayersPressedKeys - a map containing an array of player keys pressed.
var PlayersPressedKeys map[string][]string

// Init - a method that initializes service variables
func Init() {
	Base = base.Base{
		Players: []player.Player{},
		Effects: []effect.Effect{},
		Bullets: []bullet.Bullet{},
		Scene:   scene.Scene{},
		Blocks:  []block.Block{},
	}

	PlayersPressedKeys = make(map[string][]string)
}
