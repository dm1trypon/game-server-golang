package movement

import (
	"errors"
	"math"
	"time"

	"github.com/dm1trypon/game-server-golang/gameobjects"
	"github.com/ivahaev/go-logger"
)

const (
	// LC - Logging category
	LC = "[Movement] >> "
)

// ConflictControlKeys contains unallowed keys
var ConflictControlKeys = [2][2]string{{"left", "right"}, {"up", "down"}}

// IsConflictControl method checks repeated conflicts pressed keys
func IsConflictControl(keys []string) error {
	for _, checkKeys := range ConflictControlKeys {
		index := 0
		for _, key := range checkKeys {
			if len(keys) < index+1 {
				break
			}

			if key != keys[index] {
				break
			}

			index++
		}

		if index == 2 {
			textErr := "Found conflict of control player, skiped"
			logger.Warn(LC + textErr)
			return errors.New(textErr)
		}
	}

	return nil
}

// MovingPlayer method controlling speed of player
func MovingPlayer(nickname string, key string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		playerSpeed := 0

		if key == "up" {
			gameobjects.Players[index].Speed.Y--
			if !gameobjects.PressedKeys[nickname].Up {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.Y
		} else if key == "down" {
			gameobjects.Players[index].Speed.Y++
			if !gameobjects.PressedKeys[nickname].Down {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.Y
		} else if key == "left" {
			gameobjects.Players[index].Speed.X--
			if !gameobjects.PressedKeys[nickname].Left {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.X
		} else if key == "right" {
			gameobjects.Players[index].Speed.X++
			if !gameobjects.PressedKeys[nickname].Right {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.X
		} else {
			break
		}

		playerSpeedMax := gameobjects.Players[index].Speed.Max

		if int(math.Abs(float64(playerSpeed))) >= playerSpeedMax {
			break
		}

		time.Sleep(40)
	}
}

// BrakingPlayer method slows down a player
func BrakingPlayer(nickname string, key string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		playerSpeed := 0

		if key == "left" || key == "right" {
			if gameobjects.Players[index].Speed.X > 0 {
				gameobjects.Players[index].Speed.X--
			} else if gameobjects.Players[index].Speed.X < 0 {
				gameobjects.Players[index].Speed.X++
			} else if gameobjects.Players[index].Speed.X == 0 {
				return
			}

			if gameobjects.PressedKeys[nickname].Left {
				break
			}

			if gameobjects.PressedKeys[nickname].Right {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.X
		} else if key == "up" || key == "down" {
			if gameobjects.Players[index].Speed.Y > 0 {
				gameobjects.Players[index].Speed.Y--
			} else if gameobjects.Players[index].Speed.Y < 0 {
				gameobjects.Players[index].Speed.Y++
			} else if gameobjects.Players[index].Speed.Y == 0 {
				return
			}

			if gameobjects.PressedKeys[nickname].Up {
				break
			}

			if gameobjects.PressedKeys[nickname].Down {
				break
			}

			playerSpeed = gameobjects.Players[index].Speed.Y
		} else {
			return
		}

		if int(math.Abs(float64(playerSpeed))) <= 0 {
			break
		}

		time.Sleep(40)
	}
}
