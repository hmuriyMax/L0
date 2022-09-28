package main

import (
	"fmt"
)

func main() {

	for {
		var command string
		_, _ = fmt.Scanln(command)
		if command == "exit" {
			break
		}
	}
}
