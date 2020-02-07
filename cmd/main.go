package main

import (
	"os"
	"sync"

	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/dm1trypon/game-server-golang/udpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Main] >> "

func initEngine() {
	logger.Notice(LC + "Initializing GameEngine")
	var wgEngine sync.WaitGroup
	go func() {
		engine.Init()
		wgEngine.Add(1)
	}()
	wgEngine.Wait()
}

func initUDPServer() {
	logger.Notice(LC + "Initializing UDPServer")
	if err := udpserver.Start(config.GameConfig.Net.UDPPath); err != nil {
		os.Exit(1)
	}
}

func initTCPServer() {
	logger.Notice(LC + "Initializing TCPServer")
	var wgTCPServer sync.WaitGroup
	go func() {
		if err := tcpserver.Start(config.GameConfig.Net.TCPPath); err != nil {
			os.Exit(1)
		}
		wgTCPServer.Add(1)
	}()
	wgTCPServer.Wait()
}

func main() {
	logger.SetLevel("Info")
	logger.Notice(LC + "[STARTING GAME SERVER]")

	if err := config.SetConfig("config.json"); err != nil {
		logger.Error(LC + "Can not set config: " + err.Error())
		return
	}

	initEngine()
	initTCPServer()
	initUDPServer()
}
