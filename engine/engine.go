package engine

import (
	"time"

	"github.com/dm1trypon/game-server-golang/servicedata"
	"github.com/dm1trypon/game-server-golang/tcpserver"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Engine] >> "

// Start - a method that starts the main processing cycle of object timers.
func Start() {
	logger.Notice(LC + "Starting engine")

	fpsTick := 0
	secondTick := 0
	speedCalcTick := 0

	fps := servicedata.GameConfig.Game.Timers.FPS
	second := servicedata.GameConfig.Game.Timers.Second
	speedCalc := servicedata.GameConfig.Game.Timers.SpeedCalc

	// Main loop
	for {
		if fpsTick == fps {
			fpsTick = -1
		}

		if speedCalcTick == speedCalc {
			speedCalcTick = -1
		}

		if secondTick == second {
			setTimersTCPClients()
			secondTick = -1
		}

		fpsTick++
		secondTick++
		speedCalcTick++

		time.Sleep(time.Millisecond)
	}
}

func setTimersTCPClients() {
	for conn := range servicedata.TCPClients {
		servicedata.TCPClients[conn]--

		if servicedata.TCPClients[conn] == 0 {
			tcpserver.DeleteClientFromList(conn)
		}
	}
}
