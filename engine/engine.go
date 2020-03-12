package engine

import (
	"time"

	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Engine] >> "

var tickers map[string]time.Ticker

// Start - a method that starts the main processing cycle of object timers.
func Start() {
	logger.Notice(LC + "Starting engine")

	tickers = make(map[string]time.Ticker)

	fps := servicedata.GameConfig.Game.Timers.FPS
	second := servicedata.GameConfig.Game.Timers.Second
	speedCalc := servicedata.GameConfig.Game.Timers.SpeedCalc

	// Init timers
	tickers["fps"] = *time.NewTicker(time.Duration(fps) * time.Millisecond)
	tickers["speedCalc"] = *time.NewTicker(time.Duration(speedCalc) * time.Millisecond)
	tickers["second"] = *time.NewTicker(time.Duration(second) * time.Millisecond)

	for typeTimer, ticker := range tickers {
		done := make(chan bool)

		go func(ticker time.Ticker, typeTimer string) {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if typeTimer == "fps" {
						onFPS()
					} else if typeTimer == "speedCalc" {
						onSpeedCalc()
					} else if typeTimer == "second" {
						setTimersTCPClients()
					}
				}
			}
		}(ticker, typeTimer)
	}
}

func onFPS() {
	logger.Info("fps")
}

func onSpeedCalc() {
	logger.Info("speedCalc")
}

func setTimersTCPClients() {
	logger.Info(LC + "tick")
	for conn := range servicedata.TCPClients {
		servicedata.TCPClients[conn]--

		if servicedata.TCPClients[conn] == 0 {
			tcpserver.DeleteClientFromList(conn)
		}
	}
}
