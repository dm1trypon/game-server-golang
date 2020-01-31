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
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error sending:" + err.Error())
			break
		}
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error recieving:" + err.Error())
			break
		}
		fmt.Print("Message from server: " + message)
	}
}
