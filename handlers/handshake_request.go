package handlers

import (
	"fmt"
	"log"

	pb "github.com/modular-music-server/server/message"
	"github.com/modular-music-server/server/util"
	proto "google.golang.org/protobuf/proto"
)

func HandshakeRequest(client util.Client, data []byte) {
	var message pb.HandshakeRequest
	if err := proto.Unmarshal(data, &message); err != nil {
		log.Printf("Failed to unmarshal protobuf message: %v", err)
		return
	}

	fmt.Printf("Received handshake request: %v\n", message.ProtocolVersion)

	if message.ProtocolVersion != util.PROTOCOL_VERSION {
		fmt.Println("Protocol version does not match. Responding without accepting.")
		response := &pb.HandshakeResponse{
			Accepted:        false,
			RejectionReason: "Protocol version not supported",
		}
		data, err := proto.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = util.WriteMessage(client.Connection, util.MESSAGE_HANDSHAKE_RESPONSE, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	response := &pb.HandshakeResponse{
		Accepted:        true,
		ProtocolVersion: util.PROTOCOL_VERSION,
	}
	data, err := proto.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = util.WriteMessage(client.Connection, util.MESSAGE_HANDSHAKE_RESPONSE, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
