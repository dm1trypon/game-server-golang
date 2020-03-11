package servicedata

import (
	"net"

	"github.com/dm1trypon/game-server-golang/models"
)

// GameConfig - game service config
var GameConfig models.GameConfig

// TCPClients - connected clients to the TCP server with time to break the connection in case of incorrect identification.
var TCPClients map[net.Conn]int

// Init - a method that initializes service variables
func Init() {
	TCPClients = make(map[net.Conn]int)
}
