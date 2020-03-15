package physics

import (
	"strconv"
	"testing"
	"time"

	"github.com/dm1trypon/game-server-golang/models/player"
)

func TestPlayerSpeed(t *testing.T) {
	// Creating and adding new player to game's objects
	player := &player.Player{
		Nickname: "FreeMan",
		Position: player.Position{
			X:        100,
			Y:        100,
			Rotation: 0,
		},
		Size: player.Size{
			Width:  100,
			Height: 100,
		},
		Speed: player.Speed{
			X:   0,
			Y:   0,
			Max: 10,
		},
		Effects: []player.Effect{},
		Ammunition: player.Ammunition{
			Weapon: "blaster",
			Cartridges: player.Cartridges{
				Blaster: 0,
				Plazma:  0,
				Minigun: 0,
				Shotgun: 0,
			},
		},
		GameStats: player.GameStats{
			Kills: 0,
			Death: 0,
		},
		LifeStats: player.LifeStats{
			Health: 100,
			Armor:  100,
		},
	}

	keys := []string{"left"}
	speed := 0

	for i := 0; i < 100; i++ {
		PlayerControl(player, keys)

		if speed < player.Speed.Max {
			speed++
		}

		if player.Speed.X != speed {
			t.Error("Expected "+strconv.Itoa(speed)+", got ", player.Speed.X)
		}

		time.Sleep(10 * time.Millisecond)
	}

	keys = []string{}
	for i := 0; i < 10; i++ {
		PlayerControl(player, keys)

		if speed > 0 {
			speed--
		}

		if player.Speed.X != speed {
			t.Error("Expected "+strconv.Itoa(speed)+", got ", player.Speed.X)
		}
	}

	time.Sleep(100 * time.Millisecond)
}
