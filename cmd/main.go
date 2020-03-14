package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Main] >> "

func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-c
		logger.Notice(LC + "<<< STOPPING SERVICE >>>")
		os.Exit(0)
	}()
}

func main() {
	setupCloseHandler()
	logger.SetLevel("Info")
	logger.Notice(LC + "<<< STARTING SERVICE >>>")
	if !config.IsValidConfig("./config.json", "./config.schema.json") {
		return
	}

	servicedata.Init()
	engine.Start()
	tcpserver.Start()
}
