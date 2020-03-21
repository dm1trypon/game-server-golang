package manager

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"

	"github.com/dm1trypon/game-server-golang/engine"
	"github.com/dm1trypon/game-server-golang/models/client"
	"github.com/dm1trypon/game-server-golang/servicedata"

	"github.com/ivahaev/go-logger"
)

// LC - Logging category
const LC = "[Manager] >> "

// OnTCPMessage - a method that works when a message is received from
// the game client via a TCP socket, parsing and returning data.
func OnTCPMessage(msg []byte, conn net.Conn) []byte {
	connData := servicedata.GetConnData(conn)
	if connData == nil {
		return warningNotify("error", "Internal error")
	}

	time := connData.TimeDisc

	var data map[string]interface{}

	if err := json.Unmarshal(msg, &data); err != nil {
		return warningNotify("error", err.Error())
	}

	if data["method"] == nil {
		return warningNotify("error", "Method is null")
	}

	if reflect.TypeOf(data["method"]).String() != "string" {
		return warningNotify("error", "Method is not a string")
	}

	method := data["method"].(string)
	if method == "init_tcp" {
		return onTCPInit(msg, connData)
	} else if time > 0 {
		return warningNotify("error", conn.RemoteAddr().String()+" is unautorized")
	} else if method == "mouse" {
		return onMouse()
	} else if method == "keyboard" {
		return onKeyboard()
	}

	return warningNotify("error", "Method \""+method+"\" is unsupported")
}

// OnUDPMessage - a method that works when a message is received from
// the game client via a UDP socket, parsing and working with data.
func OnUDPMessage(msg []byte, udpAddr net.UDPAddr) error {
	var data map[string]interface{}

	if err := json.Unmarshal(msg, &data); err != nil {
		logger.Warn(LC + err.Error())
		return err
	}

	if data["method"] == nil {
		errText := "Method is null"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	if reflect.TypeOf(data["method"]).String() != "string" {
		errText := "Method is not a string"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	method := data["method"].(string)

	if method != "init_udp" {
		errText := "Method \"" + method + "\" is unsupported"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	if data["uuid"] == nil {
		errText := "Uuid is null"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	if reflect.TypeOf(data["uuid"]).String() != "string" {
		errText := "Uuid is not a string"
		logger.Warn(LC + errText)
		return errors.New(errText)
	}

	UUID := data["uuid"].(string)

	return engine.InitUDPClient(udpAddr, UUID)
}

func warningNotify(method string, errText string) []byte {
	logger.Warn(LC + errText)
	response, err := makeResponse(false, method, errText)
	if err != nil {
		logger.Error(LC + err.Error())
		return []byte("{\"method\":\"error\",\"message\":" + err.Error() + ",\"success\":false}")
	}

	return response
}

func onTCPInit(msg []byte, connData *servicedata.ConnectedData) []byte {
	initTCP := client.InitTCP{
		Nickname: "",
		Method:   "",
		Resolution: client.Resolution{
			Width:  -1,
			Height: -1,
		},
	}

	if err := json.Unmarshal(msg, &initTCP); err != nil {
		return warningNotify("init_tcp", err.Error())
	}

	if err := checkBodyTCPInitJSON(initTCP); err != nil {
		return warningNotify("init_tcp", err.Error())
	}

	if err := engine.InitTCPClient(initTCP, connData.TCPConnect); err != nil {
		return warningNotify("init_tcp", err.Error())
	}

	connData.TimeDisc = -1

	response, _ := makeResponse(true, "init_tcp", "OK")
	return response
}

func checkBodyTCPInitJSON(initTCP client.InitTCP) error {
	if initTCP.Nickname == "" {
		return errors.New("Nickname is empty")
	}

	if initTCP.Resolution.Width < 0 {
		return errors.New("Incorrected width")
	}

	if initTCP.Resolution.Height < 0 {
		return errors.New("Incorrected height")
	}

	return nil
}

func onMouse() []byte {
	return []byte("")
}

func onKeyboard() []byte {
	return []byte("")
}

func makeResponse(success bool, method string, message string) ([]byte, error) {
	response := client.Response{
		Method:  method,
		Message: message,
		Success: success,
	}

	msg, err := json.Marshal(response)
	if err != nil {
		logger.Error(LC + err.Error())
		return []byte(""), err
	}

	return msg, nil
}
