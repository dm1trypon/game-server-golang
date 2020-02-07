package engine

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/dm1trypon/game-server-golang/gameobjects"
	"github.com/dm1trypon/game-server-golang/models/protocol/client"
	"github.com/dm1trypon/game-server-golang/models/protocol/player"
	"github.com/dm1trypon/game-server-golang/physics/players/movement"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Engine] >> "

// Init method starts the game engine
func Init() {
	gameobjects.TDisconnect = make(map[string]*time.Timer)
	gameobjects.UDPGameClients = make(map[*client.UDPNetData]*net.UDPConn)
	gameobjects.TCPGameClients = make(map[*client.TCPNetData]*net.Conn)
	gameobjects.PressedKeys = make(map[string]*gameobjects.Control)

	gameobjects.OnInitEngine()

	logger.Notice(LC + "GameEngine has been inited")
}

// DisconnectPlayer method disconnect player by nickname
func DisconnectPlayer(addr net.Addr) error {
	isSuccess := false
	nickname := ""

	for tcpNetData := range gameobjects.TCPGameClients {
		if *tcpNetData.Addr == addr {
			delete(gameobjects.TCPGameClients, tcpNetData)
			isSuccess = true
			nickname = tcpNetData.Nickname
			break
		}
	}

	for udpNetData := range gameobjects.UDPGameClients {
		if udpNetData.Nickname == nickname {
			delete(gameobjects.UDPGameClients, udpNetData)
			isSuccess = true
			break
		}
	}

	for index, player := range gameobjects.Players {
		if player.Nickname == nickname {
			gameobjects.OnRemovePlayer(index)
			isSuccess = true
			break
		}
	}

	if _, ok := gameobjects.TDisconnect[nickname]; ok {
		gameobjects.TDisconnect[nickname].Stop()
		delete(gameobjects.TDisconnect, nickname)
		isSuccess = true
	}

	if !isSuccess {
		textErr := "Player [Nickname: " + nickname + "] already disconnected"
		logger.Warn(LC + textErr)
		return errors.New(textErr)
	}

	return nil
}

// InitTCPClient method check and set player on tcp server
func InitTCPClient(nickname string, addr net.Addr, conn *net.Conn) error {
	tcpNetData := &client.TCPNetData{
		Addr:     &addr,
		Nickname: nickname,
	}

	if _, ok := gameobjects.TCPGameClients[tcpNetData]; ok {
		textErr := "Game TCP Client[Nickname: " + nickname + ", " + addr.String() + "] already connected"
		logger.Warn(LC + textErr)
		return errors.New(textErr)
	}

	if _, err := getPlayer(nickname); err == nil {
		textErr := "Player [" + nickname + "] already exist"
		logger.Warn(LC + textErr)
		return errors.New(textErr)
	}

	logger.Notice(LC + "Connected TCP Game Client [Nickname: " + nickname + ", " + addr.String() + "]")
	gameobjects.TCPGameClients[tcpNetData] = conn

	onInitDiscTimer(nickname)
	gameobjects.OnNewPlayer(nickname)

	return nil
}

// InitUDPClient method check and set player on udp server
func InitUDPClient(nickname string, addr *net.UDPAddr, udpConn *net.UDPConn) error {
	udpNetData := &client.UDPNetData{
		Addr:     addr,
		Nickname: nickname,
	}

	if _, ok := gameobjects.UDPGameClients[udpNetData]; ok {
		textErr := "Player " + nickname + " already connected"
		logger.Warn(LC + textErr)
		return errors.New(textErr)
	}

	isAllowed := false

	for tcpNetData := range gameobjects.TCPGameClients {
		if tcpNetData.Nickname == nickname {
			isAllowed = true
		}
	}

	if !isAllowed {
		textErr := "Player " + nickname + " is unallowed"
		logger.Warn(LC + textErr)
		return errors.New(textErr)
	}

	gameobjects.UDPGameClients[udpNetData] = udpConn

	return nil
}

// OnMouse method check event of mouse from player
func OnMouse(nickname string) {
	go onUpdateTimer(nickname)
}

// KeyboardEvent method check event of keyboard from player
func KeyboardEvent(nickname string, keys []string) error {
	for tcpNetData := range gameobjects.TCPGameClients {
		if tcpNetData.Nickname == nickname {
			go onUpdateTimer(nickname)
			return keysParser(nickname, keys)
		}
	}

	textErr := "Player [Nickname: " + nickname + "] is not exist for keyboard method"
	logger.Warn(LC + textErr)
	return errors.New(textErr)
}

func getPlayer(nickname string) (player.Player, error) {
	for _, player := range gameobjects.Players {
		if player.Nickname == nickname {
			return player, nil
		}
	}

	textErr := "Player [Nickname: " + nickname + "] is not exist"
	logger.Warn(LC + textErr)

	return player.Player{}, errors.New(textErr)
}

func keysParser(nickname string, keys []string) error {
	var unallowedKeys []string

	var player player.Player
	var err error

	if player, err = getPlayer(nickname); err != nil {
		return err
	}

	for _, key := range keys {
		if key == "up" {
			gameobjects.PressedKeys[nickname].Up = true
			go movement.MovingPlayerUp(player)
			continue
		} else if key == "down" {
			gameobjects.PressedKeys[nickname].Down = true
			go movement.MovingPlayerDown(player)
			continue
		} else if key == "left" {
			gameobjects.PressedKeys[nickname].Left = true
			go movement.MovingPlayerLeft(player)
			continue
		} else if key == "right" {
			gameobjects.PressedKeys[nickname].Right = true
			go movement.MovingPlayerRight(player)
			continue
		} else if key == "1" {
			continue
		} else if key == "2" {
			continue
		} else {
			unallowedKeys = append(unallowedKeys, key)
		}
	}

	if gameobjects.PressedKeys[nickname].Up {
		gameobjects.PressedKeys[nickname].Up = false
	}

	if gameobjects.PressedKeys[nickname].Down {
		gameobjects.PressedKeys[nickname].Down = false
	}

	if gameobjects.PressedKeys[nickname].Left {
		gameobjects.PressedKeys[nickname].Left = false
	}

	if gameobjects.PressedKeys[nickname].Right {
		gameobjects.PressedKeys[nickname].Right = false
	}

	textErr := "Player [Nickname: " + nickname + ", Method: Keyboard ]: keys [" +
		strings.Join(unallowedKeys, ", ") + "] are not allowed"

	logger.Warn(LC + textErr)
	return errors.New(textErr)
}

func onInitDiscTimer(nickname string) {
	gameobjects.TDisconnect[nickname] = time.NewTimer(30 * time.Second)
	go onTimeExpired(nickname)
}

func onUpdateTimer(nickname string) {
	if gameobjects.TDisconnect[nickname].Stop() {
		gameobjects.TDisconnect[nickname] = time.NewTimer(30 * time.Second)
	}
}

func onTimeExpired(nickname string) {
	<-gameobjects.TDisconnect[nickname].C

	for tcpNetData := range gameobjects.TCPGameClients {
		if tcpNetData.Nickname == nickname {
			DisconnectPlayer(*tcpNetData.Addr)
			return
		}
	}
}

// GetGameData method compare and returns game's data
func GetGameData() []byte {
	return gameobjects.GetBase()
}

// GetTCPClients method gets connected clients via TCP
func GetTCPClients() map[*client.TCPNetData]*net.Conn {
	return gameobjects.TCPGameClients
}

// GetUDPClients method gets connected clients via UDP
func GetUDPClients() map[*client.UDPNetData]*net.UDPConn {
	return gameobjects.UDPGameClients
}
