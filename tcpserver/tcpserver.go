package tcpserver

import (
	"bytes"
	"net"
	"os"

	"github.com/ivahaev/go-logger"

	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/manager"
	"github.com/dm1trypon/game-server-golang/servicedata"
)

// LC - Logging category
const LC = "[TCPServer] >> "

// BufferSize - maximum buffer's size
const BufferSize = 1024

// Start - a method that starts a TCP server.
// Data for starting the server is taken from the game config.
func Start() {
	TCPPath := servicedata.GameConfig.Net.TCPPath
	listener, err := net.Listen("tcp", TCPPath)
	if err != nil {
		logger.Error(LC + err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	logger.Notice(LC + "TCP Server listening on " + TCPPath)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(LC + "Error accepting: " + err.Error())
			continue
		}

		addClientToList(conn)

		go handleRequest(conn)
	}
}

func addClientToList(conn net.Conn) {
	if _, ok := servicedata.TCPClients[conn]; ok {
		return
	}

	servicedata.TCPClients[conn] = 10
}

func handleRequest(conn net.Conn) {
	for {
		if _, ok := servicedata.TCPClients[conn]; !ok {
			logger.Warn(LC + "TCP " + conn.RemoteAddr().String() + " client has been disconnected")
			break
		}

		buf := make([]byte, BufferSize)
		_, err := conn.Read(buf)
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				logger.Notice(LC + "Client has been disconnected")
				engine.DeleteClientFromList(conn)
				return
			}

			logger.Error(LC + "An error occurred while receiving data from the client: " + err.Error())
			engine.DeleteClientFromList(conn)
			return
		}

		buf = bytes.Trim(buf, "\x00")
		buf = bytes.Trim(buf, "\n\t")

		logger.Info(LC + "RECV: " + string(buf))
		manager.OnTCPMessage(buf, conn)

		// conn.Write([]byte("OK\n"))
	}
}
