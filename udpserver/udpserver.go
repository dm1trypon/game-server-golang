package udpserver

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/dm1trypon/game-server-golang/servicedata"

	"github.com/dm1trypon/game-server-golang/protoworker"
	"github.com/ivahaev/go-logger"
)

const (
	// LC - Logging category
	LC = "[UDPServer] >> "
	// MAXSIZE - max size of bytes message
	MAXSIZE = 1024
	// CONNTYPE - type of server's protocol
	CONNTYPE = "udp"
)

func sender() {
	ms := 0

	for {
		if ms%17 == 0 {
			protoworker.OnFPS()
			buf := protoworker.GetGameData()
			udpGameClients := protoworker.GetUDPClients()

			for udpNetData, udpConn := range udpGameClients {
				if len(string(buf)) < 1 {
					break
				}
				// logger.Info(LC + "SENT [" + udpNetData.Addr.String() + "]: " + string(buf))
				_, err := udpConn.WriteToUDP(buf, udpNetData.Addr)
				if err != nil {
					logger.Error(LC + "Error writing message: " + err.Error())
				}
			}
		}

		if ms%100 == 0 {

		}

		if ms > 999 {
			ms = 0
		}

		ms++

		time.Sleep(1 * time.Millisecond)
	}
}

// Start method starts UDP server
func Start() {
	udpAddr, err := net.ResolveUDPAddr("udp4", servicedata.GameConfig.Net.UDPPath)
	if err != nil {
		logger.Error(LC + "Error resolving: " + err.Error())
		os.Exit(1)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Error(LC + "Error listening: " + err.Error())
		os.Exit(1)
	}
	defer udpConn.Close()

	servicedata.UDPConn = udpConn

	logger.Notice(LC + "UDP game server has been started on " + path)

	var wgSender sync.WaitGroup
	go func() {
		sender()
		wgSender.Add(1)
	}()
	wgSender.Wait()

	for {
		buf := make([]byte, MAXSIZE)

		_, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logger.Warn(LC + "Error reading message: " + err.Error())
			continue
		}

		logger.Info(LC + "RECV [" + addr.String() + "]: " + string(buf))
		logger.Info(LC + "SENT [" + addr.String() + "]: " + string(protoworker.OnUDPMessage(buf, addr, udpConn)))
	}
}
