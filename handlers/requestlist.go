package handlers

import (
	"fmt"
	"log"

	pb "github.com/modular-music-server/server/message"
	"github.com/modular-music-server/server/util"
	proto "google.golang.org/protobuf/proto"
)

func RequestList(client util.Client, data []byte) {
	var message pb.RequestList
	if err := proto.Unmarshal(data, &message); err != nil {
		log.Printf("Failed to unmarshal protobuf message: %v", err)
		return
	}

	switch message.Type {
	case pb.ListType_PROVIDERS:
		listProviders(client)
	}

	fmt.Printf("Received request list with the following type: ")
	fmt.Println(message.Type)
}

func listProviders(client util.Client) {
	var providers []*pb.ProviderEntry
	for id, provider := range client.Config.Modules.Providers {
		providers = append(providers, &pb.ProviderEntry{
			Name: provider.Name,
			Id:   id,
		})
	}

	list := &pb.ListProviders{
		Providers: providers,
	}

	data, err := proto.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = util.WriteMessage(client.Connection, util.MESSAGE_LISTPROVIDERS, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
