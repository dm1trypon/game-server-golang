package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var KEYS = []string{"left", "right", "up", "down"}

func main() {
	conn, err := net.Dial("tcp", "10.23.0.59:3001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	step := 1

	for {
		if step == 2 {
			fmt.Fprintf(conn, `{"method":"keyboard","nickname":"playernickname","keys":[`+getRandomKeys()+`]}`+"\n")
			step = 1
		}

		if step == 1 {
			fmt.Fprintf(conn, `{"method":"init_tcp","nickname":"playernickname","resolution":{"width":1280,"height":720}}`+"\n")
			step = 2
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error recieving:" + err.Error())
			break
		}
		fmt.Println("Message from server: " + message)
		time.Sleep(500 * time.Millisecond)
	}
}

func onMouse(conn net.Conn) {
	// for {
	// 	fmt.Fprintf(conn, `{"method":"mouse","nickname":"playernickname","position":{"x": `+strconv.Itoa(rand.Intn(1000))+`, "y": `+strconv.Itoa(rand.Intn(1000))+`}, "is_clicked":true}`+"\n")
	// 	message, err := bufio.NewReader(conn).ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error recieving:" + err.Error())
	// 		// break
	// 		return
	// 	}
	// 	fmt.Println("Message from server: " + message)
	// 	time.Sleep(1000 * time.Millisecond)
	// }

	fmt.Fprintf(conn, `{"method":"mouse","nickname":"playernickname","position":{"x": `+strconv.Itoa(rand.Intn(1000))+`, "y": `+strconv.Itoa(rand.Intn(1000))+`}, "is_clicked":true}`+"\n")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error recieving:" + err.Error())
		// break
		return
	}
	fmt.Println("Message from server: " + message)
	time.Sleep(1000 * time.Millisecond)
}

func getRandomKeys() string {
	var keys []string

	for step := 0; step < rand.Intn(4); step++ {
		keys = append(keys, "\""+KEYS[rand.Intn(4)]+"\"")
	}

	// fmt.Println("Random: " + strings.Join(keys, ","))

	return strings.Join(keys, ",")
}
