package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключаемся к сокету
	conn, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		// Чтение входных данных от stdin
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Text to send: ")
		// text, err := reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Println("Error sending:" + err.Error())
		// 	break
		// }
		// Отправляем в socket
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press")
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error sending:" + err.Error())
			break
		}
		fmt.Fprintf(conn, `{"method":"init_tcp","nickname":"playernickname","resolution":{"width":1280,"height":720}}`+"\n")
		// Прослушиваем ответ
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error recieving:" + err.Error())
			break
		}
		fmt.Print("Message from server: " + message)
		// fmt.Fprintf(conn, `{"method":"disconnect","nickname":"playernickname"`+"\n")
	}
}
