package gameobjects

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/models/protocol/base"
	"github.com/dm1trypon/game-server-golang/models/protocol/block"
	"github.com/dm1trypon/game-server-golang/models/protocol/bullet"
	"github.com/dm1trypon/game-server-golang/models/protocol/client"
	"github.com/dm1trypon/game-server-golang/models/protocol/effect"
	"github.com/dm1trypon/game-server-golang/models/protocol/player"
	"github.com/dm1trypon/game-server-golang/models/protocol/scene"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[GameObjects] >> "

// Control struct contains data of pressed keys
type Control struct {
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

// TDisconnect - hang timers for players
var TDisconnect map[string]*time.Timer
var PressedKeys map[string]*Control

var TCPGameClients map[*client.TCPNetData]*net.Conn
var UDPGameClients map[*client.UDPNetData]*net.UDPConn

var Players []player.Player
var Effects []effect.Effect
var Bullets []bullet.Bullet
var Blocks []block.Block
var Scene scene.Scene
var Base base.Base

// GetPlayer method gets player's object by nickname
func GetPlayer(nickname string) (player.Player, error) {
	for _, player := range Players {
		if player.Nickname == nickname {
			return player, nil
		}
	}

	textErr := "Player [Nickname: " + nickname + "] is not exist"
	logger.Warn(LC + textErr)

	return player.Player{}, errors.New(textErr)
}

// GetIndexPlayer method gets player's index nickname
func GetIndexPlayer(nickname string) (int, error) {
	for index, player := range Players {
		if player.Nickname == nickname {
			return index, nil
		}
	}

	textErr := "Player [Nickname: " + nickname + "] is not exist"
	logger.Warn(LC + textErr)

	return -1, errors.New(textErr)

}

// OnInitEngine method init struct Base
func OnInitEngine() {
	Base = base.Base{
		Players: []player.Player{},
		Effects: []effect.Effect{},
		Bullets: []bullet.Bullet{},
		Scene:   scene.Scene{},
		Blocks:  []block.Block{},
	}
}

// OnRemovePlayer method remove player in gameprotocol
func OnRemovePlayer(index int) {
	Players[index] = Players[len(Players)-1]
	Players[len(Players)-1] = *new(player.Player)
	Players = Players[:len(Players)-1]
}

// OnNewPlayer method create new player
func OnNewPlayer(nickname string) {
	weapon := config.GameConfig.GameObjects.Player.Weapon
	pCartridges := config.GameConfig.GameObjects.Player.Cartridges

	cartridges := &player.Cartridges{
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

	// Creating and adding new player to game's objects
	player := &player.Player{
		Nickname: nickname,
		Position: player.Position{
			X:        rand.Intn(config.GameConfig.GameObjects.Scene.Width),
			Y:        rand.Intn(config.GameConfig.GameObjects.Scene.Height),
			Rotation: 0,
		},
		Size: player.Size{
			Width:  config.GameConfig.GameObjects.Player.Width,
			Height: config.GameConfig.GameObjects.Player.Height,
		},
		Speed: player.Speed{
			X:   0,
			Y:   0,
			Max: config.GameConfig.GameObjects.Player.Speed,
		},
		Effects: []player.Effect{},
		Ammunition: player.Ammunition{
			Weapon:     config.GameConfig.GameObjects.Player.Weapon,
			Cartridges: *cartridges,
		},
		GameStats: player.GameStats{
			Kills: 0,
			Death: 0,
		},
		LifeStats: player.LifeStats{
			Health: config.GameConfig.GameObjects.Player.Health,
			Armor:  config.GameConfig.GameObjects.Player.Armor,
		},
	}

	Players = append(Players, *player)
	Base.Players = Players

	// Setting default keys control
	PressedKeys[nickname] = &Control{
		Up:    false,
		Down:  false,
		Left:  false,
		Right: false,
	}
}

// GetBase method compare and return Base struct
func GetBase() []byte {
	data, err := json.Marshal(&Base)
	if err != nil {
		logger.Error(LC + err.Error())
		return []byte("")
	}

	return data
}

// EditPlayer method edit player's data
func EditPlayer(player player.Player) {
	for index, nextPlayer := range Players {
		if player.Nickname == nextPlayer.Nickname {
			Players[index] = player
			return
		}
	}

	logger.Warn(LC + "Player [Nickname: " + player.Nickname + "] is not exist for update")
}
