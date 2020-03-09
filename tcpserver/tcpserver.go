package tcpserver

import (
	"fmt"
	"net"
	"os"

	"github.com/ivahaev/go-logger"

	"github.com/dm1trypon/game-server-golang/servicedata"
)

// LC - Logging category
const LC = "[TCPServer] >> "

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

		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
