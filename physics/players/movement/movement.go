package movement

import (
	"math"
	"time"

	"github.com/dm1trypon/game-server-golang/gameobjects"
	"github.com/dm1trypon/game-server-golang/models/protocol/player"
)

// MovingPlayerUp method moving an object with acceleration or braking
func MovingPlayerUp(player player.Player) {
	for {
		player.Speed.Y++
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.Y))) == player.Speed.Max {
			break
		}

		if !gameobjects.PressedKeys[player.Nickname].Up {
			go movingUpBreak(player)
			break
		}
	}
}

// MovingPlayerDown method moving an object with acceleration or braking
func MovingPlayerDown(player player.Player) {
	for {
		player.Speed.Y--
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.Y))) == player.Speed.Max {
			break
		}

		if !gameobjects.PressedKeys[player.Nickname].Down {
			go movingDownBreak(player)
			break
		}
	}
}

// MovingPlayerLeft method moving an object with acceleration or braking
func MovingPlayerLeft(player player.Player) {
	for {
		player.Speed.X++
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.X))) == player.Speed.Max {
			break
		}

		if !gameobjects.PressedKeys[player.Nickname].Left {
			go movingLeftBreak(player)
			break
		}
	}
}

// MovingPlayerRight method moving an object with acceleration or braking
func MovingPlayerRight(player player.Player) {
	for {
		player.Speed.X--
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.X))) == player.Speed.Max {
			break
		}

		if !gameobjects.PressedKeys[player.Nickname].Right {
			go movingRightBreak(player)
			break
		}
	}
}

func movingUpBreak(player player.Player) {
	for {
		player.Speed.Y--
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.Y))) == 0 {
			break
		}

		if gameobjects.PressedKeys[player.Nickname].Up {
			break
		}
	}
}

func movingDownBreak(player player.Player) {
	for {
		player.Speed.Y++
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.Y))) == 0 {
			break
		}

		if gameobjects.PressedKeys[player.Nickname].Down {
			break
		}
	}

}

func movingLeftBreak(player player.Player) {
	for {
		player.Speed.X--
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.X))) == 0 {
			break
		}

		if gameobjects.PressedKeys[player.Nickname].Left {
			break
		}
	}

}

func movingRightBreak(player player.Player) {
	for {
		player.Speed.X++
		gameobjects.EditPlayer(player)

		time.Sleep(100 * time.Millisecond)
		if int(math.Abs(float64(player.Speed.X))) == 0 {
			break
		}

		if gameobjects.PressedKeys[player.Nickname].Right {
			break
		}
	}
}
