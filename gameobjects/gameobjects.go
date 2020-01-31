package gameobjects

import (
	"github.com/dm1trypon/backend/game-server-golang/models/protocol/block"
	"github.com/dm1trypon/backend/game-server-golang/models/protocol/bullet"
	"github.com/dm1trypon/backend/game-server-golang/models/protocol/effect"
	"github.com/dm1trypon/backend/game-server-golang/models/protocol/player"
	"github.com/dm1trypon/backend/game-server-golang/models/protocol/scene"
)

var Players []player.Player
var Effects []effect.Effect
var Bullets []bullet.Bullet
var Blocks []block.Block
var Scene scene.Scene
