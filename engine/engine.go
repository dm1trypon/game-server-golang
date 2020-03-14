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

var directions = [4]string{"left", "right", "up", "down"}

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

func getMissedKeys(keys []string) []string {
	var missedKeys []string

	for _, direction := range directions {
		isExist := false

		for _, key := range keys {
			if direction == key {
				isExist = true
				break
			}
		}

		if !isExist {
			missedKeys = append(missedKeys, direction)
		}
	}

	return missedKeys
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

		isPressedHorizontal, isPressedVertical := onRacing(player, keys)
		onBraking(player, keys, isPressedVertical, isPressedHorizontal)
	}
}

func onRacing(player *player.Player, keys []string) (bool, bool) {
	isPressedVertical := false
	isPressedHorizontal := false

	for _, key := range keys {
		isPressedHorizontal, isPressedVertical =
			racing(player, key, isPressedVertical, isPressedHorizontal)
	}

	return isPressedHorizontal, isPressedVertical
}

func onBraking(player *player.Player, keys []string, isPressedVertical bool, isPressedHorizontal bool) {
	missedKeys := getMissedKeys(keys)
	for _, key := range missedKeys {
		if isBrakeVertical(key, isPressedVertical) {
			braking(player, "vertical")
		} else if isBrakeHorizonatal(key, isPressedHorizontal) {
			braking(player, "horizontal")
		}
	}
}

func racing(player *player.Player, key string, isPressedVertical bool, isPressedHorizontal bool) (bool, bool) {
	speedMax := player.Speed.Max
	speedX := int(math.Abs(float64(player.Speed.X)))
	speedY := int(math.Abs(float64(player.Speed.Y)))

	if key == "up" {
		if speedMax <= speedY {
			return isPressedHorizontal, true
		}

		player.Speed.Y++
		return isPressedHorizontal, true
	} else if key == "down" {
		if speedMax <= speedY {
			return isPressedHorizontal, true
		}

		player.Speed.Y--
		return isPressedHorizontal, true
	} else if key == "left" {
		if speedMax <= speedX {
			return true, isPressedVertical
		}

		player.Speed.X++
		return true, isPressedVertical
	} else if key == "right" {
		if speedMax <= speedX {
			return true, isPressedVertical
		}

		player.Speed.X--
		return true, isPressedVertical
	}

	logger.Warn(LC + "Unknown racing direction")
	return isPressedHorizontal, isPressedVertical
}

func braking(player *player.Player, direction string) {
	speed := 0
	isHorizontal := false
	isVertical := false

	if direction == "horizontal" {
		speed = player.Speed.X
	} else if direction == "vertical" {
		speed = player.Speed.Y
	} else {
		logger.Warn(LC + "Unknown braking direction")
		return
	}

	if speed == 0 {
		return
	}

	if speed > 0 {
		speed--
		return
	}

	if speed < 0 {
		speed++
		return
	}

	if isHorizontal {
		player.Speed.X = speed
		return
	}

	if isVertical {
		player.Speed.Y = speed
		return
	}
}

func isBrakeVertical(key string, isPressedVertical bool) bool {
	return (key == "up" || key == "down") && !isPressedVertical
}

func isBrakeHorizonatal(key string, isPressedHorizontal bool) bool {
	return (key == "left" || key == "right") && !isPressedHorizontal
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
