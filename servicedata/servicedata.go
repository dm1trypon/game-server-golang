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

// GameConfig - game service config
var GameConfig config.GameConfig

// TCPClients - connected clients to the TCP server with time to break
// the connection in case of incorrect identification.
var TCPClients map[net.Conn]int

// Base - the structure of all game objects.
var Base base.Base

// PlayersPressedKeys - a map containing an array of player keys pressed.
var PlayersPressedKeys map[string][]string

// BufPlayersPressedKeys - buffer a map containing an array of player keys pressed.
var BufPlayersPressedKeys map[string][]string

// Init - a method that initializes service variables
func Init() {
	TCPClients = make(map[net.Conn]int)

	Base = base.Base{
		Players: []player.Player{},
		Effects: []effect.Effect{},
		Bullets: []bullet.Bullet{},
		Scene:   scene.Scene{},
		Blocks:  []block.Block{},
	}

	PlayersPressedKeys = make(map[string][]string)
	BufPlayersPressedKeys = make(map[string][]string)
}
