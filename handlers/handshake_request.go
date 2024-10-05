package handlers

import (
	"fmt"
	"log"
	"net"

    "example.com/modular-music-server/util"
	pb "example.com/modular-music-server/message"
	proto "google.golang.org/protobuf/proto"
)

func HandshakeRequest(conn net.Conn, data []byte) {
    var message pb.HandshakeRequest
    if err := proto.Unmarshal(data, &message); err != nil {
        log.Printf("Failed to unmarshal protobuf message: %v", err)
        return
    }

    fmt.Printf("Received handshake request: ")
    fmt.Println(message.ProtocolVersion)

    if message.ProtocolVersion != util.PROTOCOL_VERSION {
        response := &pb.HandshakeResponse{
            Accepted: false,
            RejectionReason: "Protocol version not supported",
        }
        data, err := proto.Marshal(response)
        if err != nil {
            fmt.Println(err)
            return
        }

        err = util.WriteMessage(conn, util.MESSAGE_HANDSHAKE_RESPONSE, data)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    response := &pb.HandshakeResponse{
        Accepted: true,
        ProtocolVersion: util.PROTOCOL_VERSION,
    }
    data, err := proto.Marshal(response)
    if err != nil {
        fmt.Println(err)
        return
    }
    err = util.WriteMessage(conn, util.MESSAGE_HANDSHAKE_RESPONSE, data)
    if err != nil {
        fmt.Println(err)
        return
    }
}
