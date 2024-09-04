package common

import (
	"net"
	"strings"
)

const (
	RECEIVE_SIZE_BYTES int = 1024
)

type Protocol struct {
	socket net.Conn
}

func CreateConnection(serverAddress string) (*Protocol, error) {
	conn, err := net.Dial("tcp", serverAddress)

	if err != nil {
		return nil, err
	}

	protocol := &Protocol{
		socket: conn,
	}

	return protocol, nil
}

func (prot *Protocol) CloseConnection() error {
	return prot.socket.Close()
}

func (prot *Protocol) Send(message ClientMessage) error {
	data := []byte(message.Serialize())

	sentBytes := 0
	totalBytesToSend := len(data)

	for sentBytes < totalBytesToSend {
		bytes, err := prot.socket.Write(data[sentBytes:])
		if err != nil {
			return err
		}
		sentBytes += bytes
	}

	return nil
}

func (prot *Protocol) ReceiveAndCloseConnection() (*ServerMessage, error) {
	defer prot.socket.Close()
	buffer := make([]byte, RECEIVE_SIZE_BYTES)
	var message strings.Builder

	for {
		bytes, err := prot.socket.Read(buffer)
		if err != nil {
			return CreateServerMessage(false, "unexpected error"), err
		}

		message.Write(buffer[:bytes])

		msg := message.String()

		if strings.Contains(msg, "\n") {
			return DeserializeServerMessage(msg), nil
		}
	}
}
