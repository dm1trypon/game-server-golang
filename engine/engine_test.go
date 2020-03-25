package engine

import (
	"strconv"
	"testing"

	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/models/client"
	"github.com/dm1trypon/game-server-golang/servicedata"
)

func TestOnFPS(t *testing.T) {
	if !config.IsValidConfig("../config.json", "../config.schema.json") {
		t.Error("Config is invalid")
	}

	initTCP := client.InitTCP{
		Nickname: "test_user",
		Method:   "init_tcp",
		Resolution: client.Resolution{
			Width:  1920,
			Height: 1080,
		},
	}

	servicedata.Init()
	newPlayer(initTCP)

	for _, player := range servicedata.Base.Players {
		if player.Nickname == initTCP.Nickname {
			player.Position.X = 100
			player.Position.Y = 100
			player.Speed.X = 1
			player.Speed.Y = 1
		}
	}

	results := []int{101, 102, 103, 104, 105}

	for _, pos := range results {
		onFPS()

		for _, player := range servicedata.Base.Players {
			if player.Nickname == initTCP.Nickname {
				if player.Position.X != pos {
					t.Error("Expected "+strconv.Itoa(pos)+", got ", strconv.Itoa(player.Position.X))
				}

				if player.Position.Y != pos {
					t.Error("Expected "+strconv.Itoa(pos)+", got ", strconv.Itoa(player.Position.Y))
				}
			}
		}
	}
}
