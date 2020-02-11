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

// MovingPlayerUp method moving an object with acceleration or braking
func MovingPlayerUp(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		gameobjects.Players[index].Speed.Y++

		playerSpeedY := gameobjects.Players[index].Speed.Y
		playerSpeedMax := gameobjects.Players[index].Speed.Max

		if !gameobjects.PressedKeys[nickname].Up {
			break
		}

		if int(math.Abs(float64(playerSpeedY))) >= playerSpeedMax {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// MovingPlayerDown method moving an object with acceleration or braking
func MovingPlayerDown(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		gameobjects.Players[index].Speed.Y--

		playerSpeedY := gameobjects.Players[index].Speed.Y
		playerSpeedMax := gameobjects.Players[index].Speed.Max

		if !gameobjects.PressedKeys[nickname].Down {
			break
		}

		if int(math.Abs(float64(playerSpeedY))) >= playerSpeedMax {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// MovingPlayerLeft method moving an object with acceleration or braking
func MovingPlayerLeft(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		gameobjects.Players[index].Speed.X++

		playerSpeedX := gameobjects.Players[index].Speed.X
		playerSpeedMax := gameobjects.Players[index].Speed.Max

		if !gameobjects.PressedKeys[nickname].Left {
			break
		}

		if int(math.Abs(float64(playerSpeedX))) >= playerSpeedMax {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// MovingPlayerRight method moving an object with acceleration or braking
func MovingPlayerRight(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		gameobjects.Players[index].Speed.X--

		playerSpeedX := gameobjects.Players[index].Speed.X
		playerSpeedMax := gameobjects.Players[index].Speed.Max

		if !gameobjects.PressedKeys[nickname].Right {
			break
		}

		if int(math.Abs(float64(playerSpeedX))) >= playerSpeedMax {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// MovingUpDownBrake method slows down a player
func MovingUpDownBrake(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		if gameobjects.Players[index].Speed.Y > 0 {
			gameobjects.Players[index].Speed.Y--
		} else if gameobjects.Players[index].Speed.Y < 0 {
			gameobjects.Players[index].Speed.Y++
		} else if gameobjects.Players[index].Speed.Y == 0 {
			return
		}

		playerSpeedY := gameobjects.Players[index].Speed.Y

		if gameobjects.PressedKeys[nickname].Down {
			break
		}

		if int(math.Abs(float64(playerSpeedY))) <= 0 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// MovingLeftRightBrake method slows down a player
func MovingLeftRightBrake(nickname string) {
	for {
		var index int
		var err error

		if index, err = gameobjects.GetIndexPlayer(nickname); err != nil {
			return
		}

		if gameobjects.Players[index].Speed.X > 0 {
			gameobjects.Players[index].Speed.X--
		} else if gameobjects.Players[index].Speed.X < 0 {
			gameobjects.Players[index].Speed.X++
		} else if gameobjects.Players[index].Speed.X == 0 {
			return
		}

		playerSpeedX := gameobjects.Players[index].Speed.X

		if gameobjects.PressedKeys[nickname].Right {
			break
		}

		if int(math.Abs(float64(playerSpeedX))) <= 0 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}
