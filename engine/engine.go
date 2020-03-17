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

// TCPSchemaPath - path to TCP's schema
const TCPSchemaPath = "../schemas/tcp_data.schema.json"

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

// InitTCPClient - a method that initializes a client connected via a TCP socket.
func InitTCPClient(initTCP client.InitTCP, conn net.Conn) error {
	players := servicedata.Base.Players

	for _, player := range players {
		if player.Nickname == initTCP.Nickname {
			err := errors.New("Player " + player.Nickname + " already exists")
			return err
		}
	}

	newPlayer(initTCP)

	servicedata.TCPClients[conn] = -1
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
	servicedata.Base.Players = append(servicedata.Base.Players, player)
	mutex.Unlock()

	logger.Notice(LC + "New player has been created: {\"nickname\": \"" + player.Nickname + "\"," +
		" \"position\": [" + strconv.Itoa(player.Position.X) + ":" + strconv.Itoa(player.Position.Y) + "]}")
}

func onFPS() {
	for _, player := range servicedata.Base.Players {
		player.Position.X += player.Speed.X
		player.Position.Y += player.Speed.Y
	}

	for _, bullet := range servicedata.Base.Bullets {
		bullet.Position.X += bullet.Speed.X
		bullet.Position.Y += bullet.Speed.Y
	}
}

func onSpeedCalc() {
	playersByName := make(map[string]*player.Player)

	for _, player := range servicedata.Base.Players {
		playersByName[player.Nickname] = &player
	}

	for nickname, keys := range servicedata.PlayersPressedKeys {
		player, ok := playersByName[nickname]
		if !ok {
			continue
		}

		physics.PlayerControl(player, keys)
	}
}

// DeleteClientFromList - a method that deletes a client connected to a TCP server.
func DeleteClientFromList(conn net.Conn) {
	if _, ok := servicedata.TCPClients[conn]; !ok {
		return
	}

	delete(servicedata.TCPClients, conn)
	conn.Close()
}

func setTimersTCPClients() {
	for conn := range servicedata.TCPClients {
		servicedata.TCPClients[conn]--

		if servicedata.TCPClients[conn] == 0 {
			logger.Info(LC + conn.RemoteAddr().String() + ": Timer expired!")
			DeleteClientFromList(conn)
		}
	}
}
