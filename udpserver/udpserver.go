package udpserver

import (
	"net"
	"sync"
	"time"

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
	for {
		udpGameClients := protoworker.GetUDPClients()
		for udpNetData, udpConn := range udpGameClients {
			buf := protoworker.GetGameData()
			// logger.Info(LC + "SENT [" + udpNetData.Addr.String() + "]: " + string(buf))
			_, err := udpConn.WriteToUDP(buf, udpNetData.Addr)
			if err != nil {
				logger.Error(LC + "Error writing message: " + err.Error())
			}
		}

		protoworker.OnFPS()
		time.Sleep(17 * time.Millisecond)
	}
}

// Start method starts UDP server
func Start(path string) error {
	udpAddr, err := net.ResolveUDPAddr("udp4", path)
	if err != nil {
		logger.Error(LC + "Error resolving: " + err.Error())
		return err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Error(LC + "Error listening: " + err.Error())
		return err
	}
	defer udpConn.Close()

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

	return nil
}
