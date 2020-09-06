package main

import (
	"log"
	"os"

	"github.com/kevinnayar/aggregarian/receiver"
	"github.com/kevinnayar/aggregarian/sender"
)

func main() {
	const projectName = "aggregarian"
	ioMethod := os.Args[1]

	switch ioMethod {
	case "receive-all":
		if data, err := receiver.GetAll(projectName); err != nil {
			log.Fatal("Error in receive-all: ", err)
		} else {
			log.Printf("%+v\n", data)
		}

	case "receive-latest":
		if data, err := receiver.GetLatest(projectName); err != nil {
			log.Fatal("Error in receive-latest: ", err)
		} else {
			log.Printf("%+v\n", data)
		}

	case "send":
		sender.Start(projectName)

	default:
		log.Fatal("\nInvalid I/O method: '", ioMethod,
			"'.\nSupported I/O methods: 'receive-all', 'receive-latest', or 'send'",
		)
	}
}
