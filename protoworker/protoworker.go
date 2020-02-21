package protoworker

import (
	"bytes"
	"encoding/json"
	"net"

	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/models/protocol/client"
	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[ProtoWorker] >> "

// OnTCPMessage method parse data from client
func OnTCPMessage(data []byte, addr net.Addr, conn net.Conn) []byte {
	disconnect := &client.Disconnect{}
	initTCP := &client.InitTCP{}
	mouse := &client.Mouse{}
	keyboard := &client.Keyboard{}

	method := ""

	err := json.Unmarshal(data, &initTCP)
	if err == nil {
		method = initTCP.Method
	}

	err = json.Unmarshal(data, &disconnect)
	if err == nil {
		method = disconnect.Method
	}

	err = json.Unmarshal(data, &mouse)
	if err == nil {
		method = mouse.Method
	}

	err = json.Unmarshal(data, &keyboard)
	if err == nil {
		method = keyboard.Method
	}

	if method == "init_tcp" {
		return onInitTCPClient(initTCP.Nickname, &addr, &conn)
	} else if method == "disconnect" {
		return OnDisconnectPlayer(addr)
	} else if method == "keyboard" {
		return onKeyboard(keyboard.Nickname, keyboard.Keys)
	} else if method == "mouse" {
		return onMouse(mouse.Nickname, mouse.Position.X, mouse.Position.Y, mouse.IsClicked)
	}

	textErr := "Wrong method: " + method

	logger.Error(LC + textErr)

	return toResponse("error", textErr, false)
}

// OnUDPMessage method parse data from UDP client
func OnUDPMessage(data []byte, addr *net.UDPAddr, udpConn *net.UDPConn) []byte {
	initUDP := &client.InitUDP{}

	method := ""

	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&initUDP)
	if err == nil {
		method = initUDP.Method
	}

	if method == "init_udp" {
		return onInitUDPClient(initUDP.Nickname, addr, udpConn)
	}

	textErr := "Wrong UDP data: " + string(data)

	logger.Error(LC + textErr)

	return toResponse("error", textErr, false)
}

func onKeyboard(nickname string, keys []string) []byte {
	status := true
	message := "OK"

	if err := engine.KeyboardEvent(nickname, keys); err != nil {
		status = false
		message = err.Error()
	}

	return toResponse("keyboard", message, status)
}

func onMouse(nickname string, posX int, posY int, isClicked bool) []byte {
	status := true
	message := "OK"

	if err := engine.MouseEvent(nickname, posX, posY, isClicked); err != nil {
		status = false
		message = err.Error()
	}

	return toResponse("mouse", message, status)
}

func onInitUDPClient(nickname string, addr *net.UDPAddr, udpConn *net.UDPConn) []byte {
	status := true
	message := "OK"

	if err := engine.InitUDPClient(nickname, addr, udpConn); err != nil {
		status = false
		message = err.Error()
	}

	return toResponse("init_udp", message, status)
}

func onInitTCPClient(nickname string, addr *net.Addr, conn *net.Conn) []byte {
	status := true
	message := "OK"

	if err := engine.InitTCPClient(nickname, *addr, conn); err != nil {
		status = false
		message = err.Error()
	}

	return toResponse("init_tcp", message, status)
}

// OnDisconnectPlayer method disconnects's player
func OnDisconnectPlayer(addr net.Addr) []byte {
	status := true
	message := "OK"

	if err := engine.DisconnectPlayer(addr); err != nil {
		status = false
		message = err.Error()
	}

	return toResponse("disconnect", message, status)
}

// GetGameData method returns game's data
func GetGameData() []byte {
	return engine.GetGameData()
}

// GetUDPClients method gets connected clients via UDP
func GetUDPClients() map[*client.UDPNetData]*net.UDPConn {
	return engine.GetUDPClients()
}

func toResponse(method string, message string, status bool) []byte {
	response := &client.Response{
		Method:  method,
		Message: message,
		Status:  status,
	}

	data, err := json.Marshal(&response)
	if err != nil {
		logger.Error(LC + err.Error())
		return []byte("")
	}

	return data
}

// OnFPS method call on next tick frame
func OnFPS() {
	engine.CalcFrame()
}
