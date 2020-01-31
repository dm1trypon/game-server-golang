package main

import (
	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Main] >> "

func main() {
	logger.SetLevel("Info")
	logger.Notice(LC + "[STARTING GAME SERVER]")

	if err := config.SetConfig("config.json"); err != nil {
		logger.Error(LC + "Can not set config: " + err.Error())
		return
	}

	if err := tcpserver.Start(config.GameConfig.Net.TCPPath); err != nil {
		return
	}
}
