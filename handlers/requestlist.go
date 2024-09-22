package handlers

import (
	"fmt"
	"log"
	"net"

	pb "example.com/modular-music-server/message"
	proto "google.golang.org/protobuf/proto"
)

func RequestList(conn net.Conn, data []byte) {
    var message pb.RequestList
    if err := proto.Unmarshal(data, &message); err != nil {
        log.Printf("Failed to unmarshal protobuf message: %v", err)
        return
    }

    switch(message.Type) {
    case pb.ListType_PROVIDERS:
        listProviders()
    }

    fmt.Printf("Received request list with the following type: ")
    fmt.Println(message.Type);
}

func listProviders() {

}
