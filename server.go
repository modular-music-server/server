package main

import (
	"fmt"
	"net"

	"github.com/modular-music-server/server/config"
	"github.com/modular-music-server/server/handlers"
	"github.com/modular-music-server/server/util"
)

const Port string = "6065"

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error while loading config: %v\n", err)
		return
	}

	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:"+Port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", Port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn, config)
	}
}

var client util.Client

func handleClient(conn net.Conn, config *config.Config) {
	defer conn.Close()

	client.Connection = conn
	client.Config = config

	for {
		messageType, data, err := util.ReadMessage(conn)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF from client. Closing connection.")
				break
			}
			fmt.Println("Error: ", err)
			continue
		}

		switch messageType {
		case util.MESSAGE_HANDSHAKE_REQUEST:
			handlers.HandshakeRequest(client, data)
		case util.MESSAGE_REQUESTLIST:
			handlers.RequestList(client, data)
		case util.MESSAGE_REQUESTPROVIDER:
			handlers.RequestProvider(client, data)
		}
	}

}
