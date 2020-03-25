package engine

import (
	"errors"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/dm1trypon/game-server-golang/models/client"

	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/dm1trypon/game-server-golang/physics"
	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Engine] >> "

var tickers map[string]time.Ticker

// mutex is used to prevent the error of competitive sending messages from the web socket.
var mutex = &sync.Mutex{}

// Start - a method that starts the main processing cycle of object timers.
func Start() error {
	logger.Notice(LC + "Starting engine")
	tickers = make(map[string]time.Ticker)

	fps := servicedata.GameConfig.Game.Timers.FPS
	second := servicedata.GameConfig.Game.Timers.Second
	speedCalc := servicedata.GameConfig.Game.Timers.SpeedCalc

	// Init timers
	tickers["fps"] = *time.NewTicker(time.Duration(fps) * time.Millisecond)
	tickers["speedCalc"] = *time.NewTicker(time.Duration(speedCalc) * time.Millisecond)
	tickers["second"] = *time.NewTicker(time.Duration(second) * time.Millisecond)

	for typeTimer, ticker := range tickers {
		done := make(chan bool)

		go func(ticker time.Ticker, typeTimer string) {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if typeTimer == "fps" {
						onFPS()
					} else if typeTimer == "speedCalc" {
						onSpeedCalc()
					} else if typeTimer == "second" {
						setTimersTCPClients()
					}
				}
			}
		}(ticker, typeTimer)
	}

	return nil
}

// InitUDPClient - a method that initializes a client connected via a UDP socket.
func InitUDPClient(connAddr net.UDPAddr, UUID string) error {
	connData := servicedata.GetConnDataByUUID(UUID)
	if connData == nil {
		errText := "Can not find a connation's data by UUID, TCP connection must be first created"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	if connData.UDPAddr.IP != nil {
		errText := "UDP client with this address already connected"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	if connData.TimeDisc != -1 {
		errText := "TCP client not yet authorized"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	connData.UDPAddr = connAddr

	return nil
}

// InitTCPClient - a method that initializes a client connected via a TCP socket.
func InitTCPClient(initTCP client.InitTCP, conn net.Conn) error {
	players := servicedata.Base.Players
	nickname := ""

	for _, player := range players {
		if player.Nickname == initTCP.Nickname {
			nickname = player.Nickname
			err := errors.New("Player " + nickname + " already exists")
			return err
		}
	}

	newPlayer(initTCP)

	connData := servicedata.GetConnData(conn)
	if connData == nil {
		err := errors.New("Player " + nickname + " has't connection")
		return err
	}

	connData.TimeDisc = -1
	return nil
}

func newPlayer(initTCP client.InitTCP) {
	weapon := servicedata.GameConfig.GameObjects.Player.Weapon
	pCartridges := servicedata.GameConfig.GameObjects.Player.Cartridges

	cartridges := player.Cartridges{
		Blaster: 0,
		Plazma:  0,
		Minigun: 0,
		Shotgun: 0,
	}

	if weapon == "blaster" {
		cartridges.Blaster = pCartridges
	} else if weapon == "plazma" {
		cartridges.Plazma = pCartridges
	} else if weapon == "minigun" {
		cartridges.Minigun = pCartridges
	} else if weapon == "shotgun" {
		cartridges.Shotgun = pCartridges
	}

	player := player.Player{
		Nickname: initTCP.Nickname,
		Position: player.Position{
			X:        rand.Intn(servicedata.GameConfig.GameObjects.Scene.Width),
			Y:        rand.Intn(servicedata.GameConfig.GameObjects.Scene.Height),
			Rotation: 0,
		},
		Size: player.Size{
			Width:  servicedata.GameConfig.GameObjects.Player.Width,
			Height: servicedata.GameConfig.GameObjects.Player.Height,
		},
		Speed: player.Speed{
			X:   0,
			Y:   0,
			Max: servicedata.GameConfig.GameObjects.Player.Speed,
		},
		Effects: []player.Effect{},
		Ammunition: player.Ammunition{
			Weapon:     servicedata.GameConfig.GameObjects.Player.Weapon,
			Cartridges: cartridges,
		},
		GameStats: player.GameStats{
			Kills: 0,
			Death: 0,
		},
		LifeStats: player.LifeStats{
			Health: servicedata.GameConfig.GameObjects.Player.Health,
			Armor:  servicedata.GameConfig.GameObjects.Player.Armor,
		},
	}

	mutex.Lock()
	servicedata.Base.Players = append(servicedata.Base.Players, &player)
	mutex.Unlock()

	logger.Notice(LC + "New player has been created: {\"nickname\": \"" + player.Nickname + "\"," +
		" \"position\": [" + strconv.Itoa(player.Position.X) + ":" + strconv.Itoa(player.Position.Y) + "]}")
}

func onFPS() {
	for _, player := range servicedata.Base.Players {
		logger.Info(LC + "Position: [" + strconv.Itoa(player.Position.X) + ":" + strconv.Itoa(player.Position.Y) + "]")
		player.Position.X += player.Speed.X
		player.Position.Y += player.Speed.Y
	}

	for _, bullet := range servicedata.Base.Bullets {
		bullet.Position.X += bullet.Speed.X
		bullet.Position.Y += bullet.Speed.Y
	}

	for _, connData := range servicedata.ConnectedClients {
		if connData.UDPAddr.String() != "" {
			servicedata.UDPConn.WriteToUDP([]byte(""), &connData.UDPAddr)
		}
	}
}

func onSpeedCalc() {
	playersByName := make(map[string]*player.Player)

	for _, player := range servicedata.Base.Players {
		playersByName[player.Nickname] = player
	}

	for nickname, keys := range servicedata.PlayersPressedKeys {
		player, ok := playersByName[nickname]
		if !ok {
			continue
		}

		physics.PlayerControl(player, keys)
	}
}

func setTimersTCPClients() {
	for _, connData := range servicedata.ConnectedClients {
		connData.TimeDisc--

		if connData.TimeDisc == 0 {
			logger.Info(LC + connData.TCPConnect.RemoteAddr().String() + ": Timer expired!")
			servicedata.DelConnData(connData.TCPConnect)
		}
	}
}
