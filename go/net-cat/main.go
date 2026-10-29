package main

import (
	"NetCat/function"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var server *function.Server
	var port string

	if len(os.Args) == 2 {
		portNum, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Error: Invalid port number")
			fmt.Println("Port number must be between 1024 and 65535")
			os.Exit(1)
		}
		if portNum >= 1024 && portNum <= 65535 {
			port = os.Args[1]
		} else {
			fmt.Println("Port number must be between 1024 and 65535")
			os.Exit(1)
		}
	} else if len(os.Args) == 1 {
		port = "8989"
	} else {
		fmt.Println("Usage: go run main.go [port]")
		fmt.Println("Port number must be between 1024 and 65535")
		os.Exit(1)
	}

	server = function.NewServer(":" + port)
	fmt.Println("Listening on the port: " + port)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-server.QuitCh
}
