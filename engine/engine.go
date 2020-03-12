package engine

import (
	"math"
	"time"

	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Engine] >> "

var tickers map[string]time.Ticker

// Start - a method that starts the main processing cycle of object timers.
func Start() {
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
}

func onFPS() {
	logger.Info("fps")
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
	logger.Info("speedCalc")

	playersByName := make(map[string]*player.Player)

	for _, player := range servicedata.Base.Players {
		playersByName[player.Nickname] = &player
	}

	for nickname, keys := range servicedata.PlayersPressedKeys {
		if _, ok := playersByName[nickname]; !ok {
			continue
		}

		speedMax := playersByName[nickname].Speed.Max
		speedX := int(math.Abs(float64(playersByName[nickname].Speed.X)))
		speedY := int(math.Abs(float64(playersByName[nickname].Speed.Y)))

		for _, key := range keys {
			if key == "up" {
				if speedMax <= speedY {
					continue
				}

				playersByName[nickname].Speed.Y++
			} else if key == "down" {
				if speedMax <= speedY {
					continue
				}

				playersByName[nickname].Speed.Y--
			} else if key == "left" {
				if speedMax <= speedX {
					continue
				}

				playersByName[nickname].Speed.X++
			} else if key == "right" {
				if speedMax <= speedX {
					continue
				}

				playersByName[nickname].Speed.X--
			}
		}
	}
}

func setTimersTCPClients() {
	logger.Info(LC + "tick")
	for conn := range servicedata.TCPClients {
		servicedata.TCPClients[conn]--

		if servicedata.TCPClients[conn] == 0 {
			tcpserver.DeleteClientFromList(conn)
		}
	}
}
