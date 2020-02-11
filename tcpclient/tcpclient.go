package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var KEYS = []string{"left", "right", "up", "down"}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	step := 1

	for {
		if step == 2 {
			fmt.Fprintf(conn, `{"method":"keyboard","nickname":"playernickname","keys":["`+getRandomKeys()+`"]}`+"\n")
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

func getRandomKeys() string {
	var keys string

	for step := 0; step < rand.Intn(4); step++ {
		if step < rand.Intn(4)-1 {
			keys += KEYS[rand.Intn(4)] + ","
		} else {
			keys += KEYS[rand.Intn(4)]
		}

	}

	fmt.Println("Random: " + keys)

	return keys
}
