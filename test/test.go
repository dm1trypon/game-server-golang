package main

import (
	"fmt"
	"net"
)

func main() {

	adrr := net.UDPAddr{}

	fmt.Println(adrr.IP)
}
