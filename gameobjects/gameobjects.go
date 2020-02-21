package gameobjects

import (
	"encoding/json"
	"errors"
	"math"
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

// MoveTimers - control's timers
var MoveTimers map[string]*time.Timer

// BrakeTimers - control's timers
var BrakeTimers map[string]*time.Timer

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

// GetIndexBullet method gets bullet's index nickname
func GetIndexBullet(nickname string) (int, error) {
	for index, bullet := range Bullets {
		if bullet.Data.Nickname == nickname {
			return index, nil
		}
	}

	textErr := "Bullet [Nickname: " + nickname + "] is not exist"
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
	Players = append(Players[:index], Players[index+1:]...)
}

func getBulletSpeed(playerPosX int, playerPosY int, cursorX int, cursorY int, maxSpeed int) map[string]int {
	var speed map[string]int
	speed = make(map[string]int)
	speed["x"] = int((cursorX - playerPosX) * maxSpeed /
		int(math.Sqrt(float64((cursorX-playerPosX)*(cursorX-playerPosX)+(cursorY-playerPosY)*(cursorY-playerPosY)))))

	speed["y"] = int((cursorY - playerPosY) * maxSpeed /
		int(math.Sqrt(float64((cursorX-playerPosX)*(cursorX-playerPosX)+(cursorY-playerPosY)*(cursorY-playerPosY)))))

	return speed
}

// OnNewBullet method create player's bullet
func OnNewBullet(nickname string, cursorX int, cursorY int) {
	var index int
	var err error

	if index, err = GetIndexPlayer(nickname); err != nil {
		return
	}

	bullets := config.GameConfig.GameObjects.Bullets
	weapon := Players[index].Ammunition.Weapon
	bulletData := config.Bullet{
		Width:  -1,
		Height: -1,
		Speed:  -1,
		Weapon: "",
		Health: -1,
		Rate:   -1,
		TTL:    -1,
	}

	for _, bullet := range bullets {
		if bullet.Weapon == weapon {
			bulletData = bullet
			break
		}
	}

	if bulletData.Width == -1 {
		logger.Warn(LC + "Weapon [" + weapon + "] is not found in game's config")
		return
	}

	playerPosX := Players[index].Position.X
	playerPosY := Players[index].Position.Y

	speed := getBulletSpeed(playerPosX, playerPosY, cursorX, cursorY, bulletData.Speed)

	bullet := &bullet.Bullet{
		Position: bullet.Position{
			X:        playerPosX,
			Y:        playerPosY,
			Rotation: 0,
		},
		Size: bullet.Size{
			Width:  bulletData.Width,
			Height: bulletData.Height,
		},
		Speed: bullet.Speed{
			X:   speed["x"],
			Y:   speed["y"],
			Max: bulletData.Speed,
		},
		Data: bullet.Data{
			Nickname: nickname,
			ID:       rand.Intn(100000),
			Weapon:   weapon,
			Health:   bulletData.Health,
			TTL:      bulletData.TTL,
		},
	}

	Bullets = append(Bullets, *bullet)
	Base.Bullets = Bullets
	go OnBulletTTLExpired(nickname, bullet.Data.ID, bullet.Data.TTL)
}

// OnBulletTTLExpired method remove bullet's object when time is expired.
func OnBulletTTLExpired(nickname string, ID int, TTL int) {
	time.Sleep(time.Duration(TTL) * time.Second)

	for index, bullet := range Bullets {
		if bullet.Data.Nickname == nickname && bullet.Data.ID == ID {
			Bullets = append(Bullets[:index], Bullets[index+1:]...)
		}
	}
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
