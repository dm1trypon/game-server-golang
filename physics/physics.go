package physics

import (
	"math"

	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Physics] >> "

var directions = [4]string{"left", "right", "up", "down"}

// PlayerControl is a method that is responsible for managing a player’s object.
func PlayerControl(player *player.Player, keys []string) {
	isPressedHorizontal, isPressedVertical := onRacing(player, keys)
	onBraking(player, keys, isPressedVertical, isPressedHorizontal)
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
