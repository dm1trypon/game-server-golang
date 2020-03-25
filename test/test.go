package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	var line []string

	length := 100

	for step := 0; step < length; step++ {
		line = append(line, "*")
	}

	step := 0

	for {
		// clear()
		if step > 0 {
			line[step-1] = line[step]
			line[step] = "0"
		}

		step++

		if step > 99 {
			line[step-1] = "*"
			step = 0
		}

		fmt.Println(strings.Join(line, ""))
		time.Sleep(50 * time.Millisecond)
	}
}
