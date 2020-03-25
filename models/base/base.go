package base

import (
	"github.com/dm1trypon/game-server-golang/models/block"
	"github.com/dm1trypon/game-server-golang/models/bullet"
	"github.com/dm1trypon/game-server-golang/models/effect"
	"github.com/dm1trypon/game-server-golang/models/player"
	"github.com/dm1trypon/game-server-golang/models/scene"
)

// Base game protocol
type Base struct {
	Players []*player.Player `json:"players"`
	Bullets []*bullet.Bullet `json:"bullets"`
	Effects []*effect.Effect `json:"effects"`
	Scene   scene.Scene      `json:"scene"`
	Blocks  []*block.Block   `json:"blocks"`
}
