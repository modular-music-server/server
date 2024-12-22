package handlers

import (
	"fmt"
	"log"

	pb "github.com/modular-music-server/server/message"
	"github.com/modular-music-server/server/util"
	proto "google.golang.org/protobuf/proto"
)

func RequestProvider(client util.Client, data []byte) {
	var message pb.RequestProvider
	if err := proto.Unmarshal(data, &message); err != nil {
		log.Printf("Failed to unmarshal protobuf message: %v", err)
		return
	}

	fmt.Printf("Received provider request: ")
	fmt.Println(message.Id)

	provider, ok := client.Config.Modules.Providers[message.Id]
	if !ok {
		fmt.Println("Received provider request for provider we do not have!")
	}
	fmt.Println("We have that provider!")
	fmt.Println(provider)

}
