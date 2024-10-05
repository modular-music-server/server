package handlers

import (
	"fmt"
	"log"
	"net"

	"example.com/modular-music-server/util"
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
        listProviders(conn)
    }

    fmt.Printf("Received request list with the following type: ")
    fmt.Println(message.Type);
}

func listProviders(conn net.Conn) {
    list := &pb.ListProviders{
        Providers: []*pb.ProviderEntry{
            { Id: "youtube", Name: "YouTube" },
            { Id: "test", Name: "Test" },
        },
    }

    data, err := proto.Marshal(list)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = util.WriteMessage(conn, util.MESSAGE_LISTPROVIDERS, data)
    if err != nil {
        fmt.Println(err)
        return
    }
}
