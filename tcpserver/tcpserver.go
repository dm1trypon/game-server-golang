package tcpserver

import (
	"bufio"
	"net"
	"sync"

	"github.com/dm1trypon/game-server-golang/protoworker"
	"github.com/ivahaev/go-logger"
)

const (
	// LC - Logging category
	LC = "[TCPServer] >> "
	// CONNTYPE - type of server's protocol
	CONNTYPE = "tcp"
)

// Start method starts TCP server
func Start(path string) error {
	listener, err := net.Listen(CONNTYPE, path)
	if err != nil {
		logger.Error(LC + "Error listening: " + err.Error())
		return err
	}
	defer listener.Close()

	logger.Notice(LC + "TCP game server has been started on " + path)

	loop(listener)

	return nil
}

func loop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(LC + "Error accepting game client: " + err.Error())
			return
		}

		var wgReq sync.WaitGroup
		go func() {
			handleRequest(conn)
			wgReq.Add(1)
		}()
		wgReq.Wait()
	}
}

func handleRequest(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Error(LC + err.Error())
			onDisconnected(conn)
			break
		}

		addr := conn.RemoteAddr()
		logger.Info(LC + "RECV [" + addr.String() + "]: " + message)
		data := string(protoworker.OnTCPMessage([]byte(message), addr, conn)) + "\n"
		conn.Write([]byte(data))
	}
}

func onDisconnected(conn net.Conn) {
	conn.RemoteAddr()
	data := string(protoworker.OnDisconnectPlayer(conn.RemoteAddr())) + "\n"
	conn.Write([]byte(data))
}
