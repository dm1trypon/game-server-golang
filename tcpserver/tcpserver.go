package tcpserver

import (
	"net"

	"github.com/ivahaev/go-logger"
)

const (
	// LC - Logging category
	LC = "[TCPServer] >> "
	// MAXSIZE - max size of bytes message
	MAXSIZE = 1024
	// CONNTYPE - type of server's protocol
	CONNTYPE = "tcp"
)

var tcpGameClients map[string]net.Conn

// Start method starts TCP server
func Start(path string) error {
	tcpGameClients = make(map[string]net.Conn)

	l, err := net.Listen(CONNTYPE, path)
	if err != nil {
		logger.Error(LC + "Error listening: " + err.Error())
		return err
	}
	defer l.Close()

	logger.Notice(LC + "TCP game server has been started on " + path)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error(LC + "Error accepting game client: " + err.Error())
			return err
		}

		host := conn.RemoteAddr().String()

		if _, ok := tcpGameClients[host]; ok != false {
			logger.Warn(LC + "Game client " + host + " already connected")
			continue
		}

		logger.Notice(LC + "Connected game client: " + host)
		tcpGameClients[host] = conn

		go handleRequest(conn)
	}

	return nil
}

func handleRequest(conn net.Conn) {
	for {
		buf := make([]byte, MAXSIZE)
		_, err := conn.Read(buf)
		if err != nil {
			onDisconnected(conn)
			break
		}

		logger.Info(LC + "Incomming message from host " + conn.RemoteAddr().String() + ": " + string(buf))
		conn.Write(buf)
	}
}

func onDisconnected(conn net.Conn) {
	host := conn.RemoteAddr().String()
	if _, ok := tcpGameClients[host]; ok != true {
		logger.Warn(LC + "Game client " + host + " already disconnected")
		return
	}

	conn.Close()
	delete(tcpGameClients, host)
	logger.Notice(LC + "Disconnected game client: " + host)
}
