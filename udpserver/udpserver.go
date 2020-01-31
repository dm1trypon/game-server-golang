package main

import (
	"net"
	"time"

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

var udpGameClients map[*net.UDPAddr]*net.UDPConn

func sender() {
	buf := []byte("Test message")
	for {
		for host, udpConn := range udpGameClients {
			logger.Info(LC + "Sending message to game client " + host.String() + ": " + string(buf))
			_, err := udpConn.WriteToUDP(buf, host)
			if err != nil {
				logger.Error(LC + "Error writing message: " + err.Error())
			}
		}

		time.Sleep(1 * time.Second)
	}

}

func main() {
	logger.SetLevel("Info")

	udpGameClients = make(map[*net.UDPAddr]*net.UDPConn)

	udpAddr, err := net.ResolveUDPAddr("udp4", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		logger.Error(LC + "Error resolving: " + err.Error())
		return
	}

	logger.Notice(LC + "UDP game server has been resolved on " + CONN_HOST + ":" + CONN_PORT)

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Error(LC + "Error listening: " + err.Error())
		return
	}
	defer udpConn.Close()

	logger.Notice(LC + "UDP game server has been started on " + CONN_HOST + ":" + CONN_PORT)

	go sender()

	for {
		buf := make([]byte, MAXSIZE)
		n, host, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logger.Warn(LC + "Error reading message: " + err.Error())
			continue
		}

		logger.Info(LC + "Incomming message from game client: " + string(buf[:n]))

		if _, ok := udpGameClients[host]; ok != true {
			udpGameClients[host] = udpConn
		}

		_, err = udpConn.WriteToUDP(buf, host)
		if err != nil {
			logger.Error(LC + "Error writing message: " + err.Error())
		}
	}

}
