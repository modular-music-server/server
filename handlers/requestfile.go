package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	pb "github.com/modular-music-server/protobufs/go"
	"github.com/modular-music-server/server/util"
	proto "google.golang.org/protobuf/proto"
)

func RequestFile(client util.Client, data []byte) {
	var message pb.RequestFile;
	if err := proto.Unmarshal(data, &message); err != nil {
		log.Printf("Failed to unmarshal protobuf message: %v", err)
		return
	}

	fmt.Printf("Received file request: ")
	path, err := getFilePath(client, &message)
	if err != nil {
		log.Printf("Failed to get path of file for file request: %v", err)
		return
	}
	data, err = os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file! %s %v", path, err)
	}

	fileInfoMessage := &pb.FileInfo{
		Type: message.Type,
		Name: message.Name,
		Size: int32(len(data)),
	}

	fileInfoMessageData, err := proto.Marshal(fileInfoMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.WriteMessage(util.MESSAGE_FILEINFO, fileInfoMessageData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileChunkMessage := &pb.FileChunk{
		Type: message.Type,
		Name: message.Name,
		Size: int32(len(data)),
		Data: data,
	}

	fileChunkMessageData, err := proto.Marshal(fileChunkMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.WriteMessage(util.MESSAGE_FILECHUNK, fileChunkMessageData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getFilePath(client util.Client, message *pb.RequestFile) (string, error) {
	switch message.Type {
	case pb.FileType_PROVIDER_CLIENT:
		provider, ok := client.Config.Modules.Providers[message.Name]
		if !ok {
			return "", errors.New("Received provider client file request for provider we do not have!")
		}
		path := path.Join(provider.Location, "client.lua")
		return path, nil
	default:
		return "", errors.New("We don't have a handler for that file type!")
	}
}
