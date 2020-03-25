package udpserver

import (
	"net"
	"os"

	"github.com/dm1trypon/game-server-golang/manager"
	"github.com/dm1trypon/game-server-golang/servicedata"

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

	servicedata.UDPConn = *udpConn
	logger.Notice(LC + "UDP game server has been started on " + servicedata.GameConfig.Net.UDPPath)

	listenHandler(udpConn)
}

func listenHandler(udpConn *net.UDPConn) {
	for {
		buf := make([]byte, MAXSIZE)

		_, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logger.Warn(LC + "Error reading message: " + err.Error())
			continue
		}

		logger.Info(LC + "RECV [" + addr.String() + "]: " + string(buf))
		manager.OnUDPMessage(buf, *addr)
	}
}
