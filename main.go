package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevinnayar/aggregarian/receiver"
	"github.com/kevinnayar/aggregarian/sender"
)

func main() {
	const projectName = "aggregarian"
	ioMethod := os.Args[1]

	switch ioMethod {
	case "receive":
		fmt.Printf("Starting I/O method... '%s'\n", ioMethod)
		receiver.Start(projectName)

	case "send":
		fmt.Printf("Starting I/O method... '%s'\n", ioMethod)
		sender.Start(projectName)

	default:
		log.Fatal("\nInvalid I/O method: '", ioMethod, "'.\nSupported I/O methods: 'receive' or 'send'")
	}
}
