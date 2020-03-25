package servicedata

import (
	"net"

	"github.com/ivahaev/go-logger"

	"github.com/delaemon/go-uuidv4"
	"github.com/dm1trypon/game-server-golang/models/base"
	"github.com/dm1trypon/game-server-golang/models/block"
	"github.com/dm1trypon/game-server-golang/models/bullet"
	"github.com/dm1trypon/game-server-golang/models/config"
	"github.com/dm1trypon/game-server-golang/models/effect"
	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/dm1trypon/game-server-golang/models/scene"
)

// LC - Logging category
const LC = "[ServiceData] >> "

// ConnectedData - contains data on connections to TCP and UDP
// servers and the time to disconnect.
type ConnectedData struct {
	TCPConnect net.Conn
	UDPAddr    net.UDPAddr
	UUID       string
	TimeDisc   int
	Nickname   string
}

// GameConfig - game service config
var GameConfig config.GameConfig

// ConnectedClients - connected clients to the TCP and UDP servers with time to break
// the connection in case of incorrect identification.
var ConnectedClients []*ConnectedData

// UDPConn - UDP's connection
var UDPConn net.UDPConn

// UDPClients - connected client's addresses to the UDP server.
var UDPClients map[net.Conn]net.UDPAddr

// Base - the structure of all game objects.
var Base base.Base

// PlayersPressedKeys - a map containing an array of player keys pressed.
var PlayersPressedKeys map[string][]string

// Init - a method that initializes service variables
func Init() {
	Base = base.Base{
		Players: []*player.Player{},
		Effects: []*effect.Effect{},
		Bullets: []*bullet.Bullet{},
		Scene:   scene.Scene{},
		Blocks:  []*block.Block{},
	}

	PlayersPressedKeys = make(map[string][]string)
}

// DelConnData - a method that deletes a client connected to a TCP server.
func DelConnData(conn net.Conn) {
	var connClients []*ConnectedData
	for _, connData := range ConnectedClients {
		if connData.TCPConnect != conn {
			connClients = append(connClients, connData)
		}
	}

	ConnectedClients = connClients
	conn.Close()
}

// IsExistConnData is a method that checks for the presence of data from a connected client.
func IsExistConnData(conn net.Conn) bool {
	for _, connData := range ConnectedClients {
		if connData.TCPConnect == conn {
			return true
		}
	}

	return false
}

// AddConnData is a method that adds data to a connected client.
func AddConnData(conn net.Conn) {
	if IsExistConnData(conn) {
		return
	}

	uuid, err := uuidv4.Generate()
	if err != nil {
		logger.Error(LC + err.Error())
		return
	}

	connectedData := &ConnectedData{
		TCPConnect: conn,
		UDPAddr:    net.UDPAddr{},
		TimeDisc:   10,
		Nickname:   "",
		UUID:       uuid,
	}

	ConnectedClients = append(ConnectedClients, connectedData)
}

// GetConnData - get connection client's data
func GetConnData(conn net.Conn) *ConnectedData {
	for _, connData := range ConnectedClients {
		if connData.TCPConnect == conn {
			return connData
		}
	}

	return nil
}

// GetConnDataByUUID - get connection client's data by UUID
func GetConnDataByUUID(UUID string) *ConnectedData {
	for _, connData := range ConnectedClients {
		if connData.UUID == UUID {
			return connData
		}
	}

	return nil
}
