package engine

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/dm1trypon/game-server-golang/gameobjects"
	"github.com/dm1trypon/game-server-golang/models/protocol/client"
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

	if _, err := gameobjects.GetPlayer(nickname); err == nil {
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

func keysParser(nickname string, keys []string) error {
	log.Println(len(keys))
	if len(keys) > 2 {
		textErr := "Player [Nickname: " + nickname + ", Method: Keyboard ]: key limit exceeded"

		logger.Warn(LC + textErr)

		return errors.New(textErr)
	}

	if err := movement.IsConflictControl(keys); err != nil {
		keys = []string{}
	}

	var unallowedKeys []string

	// Used to identify the keys pressed last time.
	control := &gameobjects.Control{
		Up:    false,
		Down:  false,
		Left:  false,
		Right: false,
	}

	for _, key := range keys {
		if key == "up" {
			if gameobjects.PressedKeys[nickname].Up {
				control.Up = true
				control.Down = false
				continue
			}

			gameobjects.PressedKeys[nickname].Up = true
			gameobjects.PressedKeys[nickname].Down = false
			control.Up = true
			control.Down = false
			go movement.MovingPlayerUp(nickname)
			continue
		} else if key == "down" {
			if gameobjects.PressedKeys[nickname].Down {
				control.Down = true
				control.Up = false
				continue
			}

			gameobjects.PressedKeys[nickname].Down = true
			gameobjects.PressedKeys[nickname].Up = false
			control.Down = true
			control.Up = false
			go movement.MovingPlayerDown(nickname)
			continue
		} else if key == "left" {
			if gameobjects.PressedKeys[nickname].Left {
				control.Left = true
				control.Right = false
				continue
			}

			gameobjects.PressedKeys[nickname].Left = true
			gameobjects.PressedKeys[nickname].Right = false
			control.Left = true
			control.Right = false
			go movement.MovingPlayerLeft(nickname)
			continue
		} else if key == "right" {
			if gameobjects.PressedKeys[nickname].Right {
				control.Right = true
				control.Left = false
				continue
			}

			gameobjects.PressedKeys[nickname].Right = true
			gameobjects.PressedKeys[nickname].Left = false
			control.Right = true
			control.Left = false
			go movement.MovingPlayerRight(nickname)
			continue
		} else if key == "1" {
			continue
		} else if key == "2" {
			continue
		} else {
			unallowedKeys = append(unallowedKeys, key)
		}
	}

	if len(unallowedKeys) > 0 {
		textErr := "Player [Nickname: " + nickname + ", Method: Keyboard ]: keys [" +
			strings.Join(unallowedKeys, ", ") + "] are not allowed"

		logger.Warn(LC + textErr)

		return errors.New(textErr)
	}

	// If the key states do not match, the Braking method starts.
	if gameobjects.PressedKeys[nickname].Up && !control.Up {
		gameobjects.PressedKeys[nickname].Up = false
		go movement.MovingUpDownBrake(nickname)
	}

	if gameobjects.PressedKeys[nickname].Down && !control.Down {
		gameobjects.PressedKeys[nickname].Down = false
		go movement.MovingUpDownBrake(nickname)
	}

	if gameobjects.PressedKeys[nickname].Left && !control.Left {
		gameobjects.PressedKeys[nickname].Left = false
		go movement.MovingLeftRightBrake(nickname)
	}

	if gameobjects.PressedKeys[nickname].Right && !control.Right {
		gameobjects.PressedKeys[nickname].Right = false
		go movement.MovingLeftRightBrake(nickname)
	}

	return nil
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

// FPS method calc position of game's objects
func FPS() {
	for index := range gameobjects.Players {
		gameobjects.Players[index].Position.X += gameobjects.Players[index].Speed.X
		gameobjects.Players[index].Position.Y += gameobjects.Players[index].Speed.Y
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
