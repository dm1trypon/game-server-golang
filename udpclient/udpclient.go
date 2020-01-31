package main

import (
	"net"

	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[TCPServer] >> "

// MAXSIZE - max size of byte's array
const MAXSIZE = 1024

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		logger.Error(LC + "Error resolving: " + err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Error(LC + "Can not connect to UDP game server: " + err.Error())
		return
	}
	defer conn.Close()

	message := []byte("Connect")
	_, err = conn.Write(message)
	if err != nil {
		logger.Error(LC + "Error writing message: " + err.Error())
		return
	}

	for {
		buf := make([]byte, MAXSIZE)
		n, host, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Warn(LC + "Error reading message: " + err.Error())
			continue
		}

		logger.Info(LC + "Incomming message from game server " + host.String() + ": " + string(buf[:n]))
	}
}
