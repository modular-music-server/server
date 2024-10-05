package handlers

import (
	"fmt"
	"log"
	"net"

    // "example.com/modular-music-server/util"
	pb "example.com/modular-music-server/message"
	proto "google.golang.org/protobuf/proto"
)

func RequestProvider(conn net.Conn, data []byte) {
    var message pb.RequestProvider
    if err := proto.Unmarshal(data, &message); err != nil {
        log.Printf("Failed to unmarshal protobuf message: %v", err)
        return
    }

    fmt.Printf("Received provider request: ")
    fmt.Println(message.Id)
}
