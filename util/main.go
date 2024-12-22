package util

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/modular-music-server/server/config"
)

const PROTOCOL_VERSION = "0.0.1"

type MessageType uint8

const (
	MESSAGE_HANDSHAKE_REQUEST MessageType = iota
	MESSAGE_HANDSHAKE_RESPONSE
	MESSAGE_REQUESTLIST
	MESSAGE_LISTPROVIDERS
	MESSAGE_REQUESTPROVIDER
)

type Client struct {
	Connection net.Conn
	Config     *config.Config
}

func ReadMessage(conn net.Conn) (MessageType, []byte, error) {
	messageTypeBuf := make([]byte, 1)
	_, err := io.ReadFull(conn, messageTypeBuf)
	if err != nil {
		return 0, nil, err
	}

	sizeBuf := make([]byte, 4)
	_, err = io.ReadFull(conn, sizeBuf)
	if err != nil {
		return 0, nil, err
	}

	messageSize := binary.BigEndian.Uint32(sizeBuf)
	messageBuf := make([]byte, messageSize)

	_, err = io.ReadFull(conn, messageBuf)
	if err != nil {
		return 0, nil, err
	}

	return MessageType(messageTypeBuf[0]), messageBuf, nil
}

func WriteMessage(conn net.Conn, messageType MessageType, message []byte) error {
	sizeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuf, uint32(len(message)))

	_, err := conn.Write([]byte{byte(messageType)})
	if err != nil {
		return err
	}
	_, err = conn.Write(sizeBuf)
	if err != nil {
		return err
	}
	_, err = conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}
