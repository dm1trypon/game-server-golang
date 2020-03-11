package main

import (
	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Main] >> "

func main() {
	logger.SetLevel("Info")
	logger.Notice(LC + "<<< STARTING SERVICE >>>")
	if !config.IsValidConfig("./config.json", "./config.schema.json") {
		logger.Notice(LC + "<<< STOPING SERVICE >>>")
		return
	}

	servicedata.Init()

	go tcpserver.Start()
	engine.Start()
}
