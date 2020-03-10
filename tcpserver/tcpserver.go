package tcpserver

import (
	"io"
	"net"
	"os"

	"github.com/ivahaev/go-logger"

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

		addClientToList(conn.RemoteAddr().String())

		go handleRequest(conn)
	}
}

func addClientToList(host string) {
	if _, ok := servicedata.TCPClients[host]; ok {
		return
	}

	servicedata.TCPClients[host] = 10
}

func deleteClientFromList(host string) {
	if _, ok := servicedata.TCPClients[host]; !ok {
		return
	}

	delete(servicedata.TCPClients, host)
}

func handleRequest(conn net.Conn) {
	host := conn.RemoteAddr().String()
	buf := make([]byte, BufferSize)
	_, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			logger.Notice(LC + "Client has been disconnected")
			deleteClientFromList(host)
			conn.Close()
			return
		}

		logger.Error(LC + "An error occurred while receiving data from the client: " + err.Error())
		deleteClientFromList(host)
		conn.Close()
		return
	}

	conn.Write([]byte("OK"))
}
