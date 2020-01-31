package engine

import (
	"github.com/dm1trypon/game-server-golang/gamedata"
)

// Start method starts the game engine
func Start() {

}

func onStep() {
	gamedata.OnSend()
}
