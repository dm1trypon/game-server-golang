package servicedata

import (
	"github.com/dm1trypon/game-server-golang/models"
)

// GameConfig - game service config
var GameConfig models.GameConfig

// TCPClients - connected clients to the TCP server with time to break the connection in case of incorrect identification.
var TCPClients map[string]int
