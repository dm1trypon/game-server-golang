package manager

import (
	"bytes"
	"net"
	"testing"

	"github.com/dm1trypon/game-server-golang/config"
	"github.com/dm1trypon/game-server-golang/servicedata"
)

type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

type mockConn struct {
	net.Conn
	b *bytes.Buffer
}

type JSON struct {
	input  string
	output string
}

func TestOnTCPMessage(t *testing.T) {
	net.Listen("tcp", "127.0.0.1:3333")
	conn, _ := net.Dial("tcp", "127.0.0.1:3333")

	results := []JSON{
		JSON{
			input:  "{\"method\":\"keyboard\",\"nickname\":\"FreeMan\",\"keys\":[\"left\", \"up\"]}",
			output: "{\"method\":\"error\",\"message\":\"127.0.0.1:3333 is unautorized\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":123,\"nickname\":\"FreeMan\",\"keys\":[\"left\", \"up\"]}",
			output: "{\"method\":\"error\",\"message\":\"Method is not a string\",\"success\":false}",
		},
		JSON{
			input:  "{\"nickname\":\"FreeMan\",\"keys\":[\"left\", \"up\"]}",
			output: "{\"method\":\"error\",\"message\":\"Method is null\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":\"keyboard\",\"nickname\":\"FreeMan\",\"keys\":[\"left\", \"up\"]",
			output: "{\"method\":\"error\",\"message\":\"unexpected end of JSON input\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":\"init_tcp\",\"nickname\":\"FreeMan\",\"resolution\":{\"width\":-1,\"height\":1080}}",
			output: "{\"method\":\"init_tcp\",\"message\":\"Incorrected width\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":\"init_tcp\",\"nickname\":\"FreeMan\",\"resolution\":{\"width\":1920,\"height\":-1}}",
			output: "{\"method\":\"init_tcp\",\"message\":\"Incorrected height\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":\"init_tcp\",\"nickname\":\"\",\"resolution\":{\"width\":1920,\"height\":-1}}",
			output: "{\"method\":\"init_tcp\",\"message\":\"Nickname is empty\",\"success\":false}",
		},
		JSON{
			input:  "{\"method\":\"init_tcp\",\"nickname\":\"FreeMan\",\"resolution\":{\"width\":1920,\"height\":1080}}",
			output: "{\"method\":\"init_tcp\",\"message\":\"OK\",\"success\":true}",
		},
		JSON{
			input:  "{\"method\":\"hello\",\"nickname\":\"FreeMan\",\"keys\":[\"left\", \"up\"]}",
			output: `{"method":"error","message":"Method \"hello\" is unsupported","success":false}`,
		},
		JSON{
			input:  "{\"method\":\"init_tcp\",\"nickname\":\"FreeMan\",\"resolution\":{\"width\":1920,\"height\":1080}}",
			output: "{\"method\":\"init_tcp\",\"message\":\"Player FreeMan already exists\",\"success\":false}",
		},
	}

	if !config.IsValidConfig("../config.json", "../config.schema.json") {
		t.Error("Config is invalid")
	}

	servicedata.Init()
	servicedata.TCPClients[conn] = 10

	for _, result := range results {
		msg := string(OnTCPMessage([]byte(result.input), conn))
		if result.output != msg {
			t.Error("Expected "+result.output+", got ", msg)
		}
	}
}
