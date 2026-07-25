package main

import (
	"fmt"
	"os"

	"github.com/saltlakecity/gostart/internal/commands"
)

func main() {
	command := os.Args[1]

	switch command {
	case "create":
		commands.RunCreate()
	default:
		fmt.Printf(`
		Invalid command: %s
		
		To create a project try

		<gostart create>
		`, command)
	}
}
