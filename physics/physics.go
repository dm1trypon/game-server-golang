package physics

import (
	"math"

	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Physics] >> "

var orientations = [4]string{"horizontal", "vertical"}
var hDirections = [2]string{"left", "right"}
var vDirections = [2]string{"up", "down"}

// PlayerControl is a method that is responsible for managing a player’s object.
func PlayerControl(player *player.Player, keys []string) {
	isPressedHorizontal, isPressedVertical := onRacing(player, keys)
	onBraking(player, keys, isPressedVertical, isPressedHorizontal)
}

func isBrakingOrientation(orientation string, keys []string) bool {
	var directions [2]string

	if orientation == "horizontal" {
		directions = hDirections
	} else if orientation == "vertical" {
		directions = vDirections
	} else {
		logger.Warn("Unknown braking direction")
		return false
	}

	for _, direction := range directions {
		isExist := false

		for _, key := range keys {
			if direction == key {
				isExist = true
				break
			}
		}

		if !isExist {
			return true
		}

		break
	}

	return false
}

func getBrakingOrientations(keys []string) []string {
	var brakingOrientations []string

	for _, orientation := range orientations {
		if isBrakingOrientation(orientation, keys) {
			brakingOrientations = append(brakingOrientations, orientation)
		}
	}

	return brakingOrientations
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
	brakingOrientations := getBrakingOrientations(keys)

	for _, orientation := range brakingOrientations {
		if !isPressedVertical {
			braking(player, orientation)
		} else if !isPressedHorizontal {
			braking(player, orientation)
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

func braking(player *player.Player, orientation string) {
	speed := 0
	isHorizontal := false
	isVertical := false

	if orientation == "horizontal" {
		speed = player.Speed.X
		isHorizontal = true
	} else if orientation == "vertical" {
		speed = player.Speed.Y
		isVertical = true
	} else {
		logger.Warn(LC + "Unknown braking direction")
		return
	}

	if speed == 0 {
		return
	}

	if speed > 0 {
		speed--
	}

	if speed < 0 {
		speed++
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
