package base

import (
	"github.com/dm1trypon/game-server-golang/models/protocol/block"
	"github.com/dm1trypon/game-server-golang/models/protocol/bullet"
	"github.com/dm1trypon/game-server-golang/models/protocol/effect"
	"github.com/dm1trypon/game-server-golang/models/protocol/player"
	"github.com/dm1trypon/game-server-golang/models/protocol/scene"
)

// Base game protocol
type Base struct {
	Players []player.Player `json:"players"`
	Bullets []bullet.Bullet `json:"bullets"`
	Effects []effect.Effect `json:"effects"`
	Scene   scene.Scene     `json:"scene"`
	Blocks  []block.Block   `json:"blocks"`
}
