package main

import (
	"github.com/dm1trypon/game-server-golang/config"
	"github.com/ivahaev/go-logger"
)

func main() {
	logger.SetLevel("Info")
	logger.Notice("<<< STARTING SERVICE >>>")
	if !config.IsValidConfig("../config.json", "../config.schema.json") {
		logger.Notice("<<< STOPING SERVICE >>>")
		return
	}
}
