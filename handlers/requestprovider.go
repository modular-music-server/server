package handlers

import (
	"fmt"
	"log"

	pb "example.com/modular-music-server/message"
	"example.com/modular-music-server/util"
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
}
