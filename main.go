package main

import (
	"log"
	"os"
)

func main() {
	cliArgs := os.Args[1:]
	if len(cliArgs) == 0 {
		log.Fatalln("No arguments provided")
	}
	if cliArgs[0] == "server" {
		runTCPServer()
	} else if cliArgs[0] == "client" {
		runTCPClient(false)
	} else if cliArgs[0] == "admin" {
		runTCPClient(true)
	} else {
		log.Fatalln("Invalid argument")
	}
}
