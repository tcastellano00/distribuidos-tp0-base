package common

import (
	"bytes"
	"net"
	"strings"
)

const (
	RECEIVE_SIZE_BYTES int = 1024
	MAX_PAYLOAD_SIZE   int = 8 * 1024 //8kb
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

func (prot *Protocol) SendBetList(messages []*ClientMessage) error {
	var buffer bytes.Buffer

	//Message type
	buffer.Write([]byte("bet\n"))

	//Serializo todos los mensajes
	for _, message := range messages {
		serialized := message.Serialize()
		data := []byte(serialized)

		_, err := buffer.Write(data)
		if err != nil {
			return err
		}
	}

	//End of message
	buffer.Write([]byte("\n\n"))

	payload := buffer.Bytes()

	//Divido en chunks
	for start := 0; start < len(payload); start += MAX_PAYLOAD_SIZE {
		end := start + MAX_PAYLOAD_SIZE
		if end > len(payload) {
			end = len(payload)
		}

		// Obtiene el fragmento de datos
		chunk := payload[start:end]

		// Envía el fragmento
		if err := prot.sendAll(chunk); err != nil {
			return err
		}
	}

	return nil
}

func (prot *Protocol) SendBet(message ClientMessage) error {
	var buffer bytes.Buffer
	//Message type
	buffer.Write([]byte("bet\n"))
	buffer.Write([]byte(message.Serialize()))
	buffer.Write([]byte("\n\n"))

	payload := buffer.Bytes()

	return prot.sendAll(payload)
}

func (prot *Protocol) BetReady(clientId string) error {
	var buffer bytes.Buffer
	//Message type
	buffer.Write([]byte("ready\n"))
	buffer.Write([]byte(clientId + "\n"))
	buffer.Write([]byte("\n\n"))

	payload := buffer.Bytes()

	return prot.sendAll(payload)
}

func (prot *Protocol) AskForResults(clientId string) error {
	var buffer bytes.Buffer
	//Message type
	buffer.Write([]byte("results\n"))
	buffer.Write([]byte(clientId + "\n"))
	buffer.Write([]byte("\n\n"))

	payload := buffer.Bytes()

	return prot.sendAll(payload)
}

func (prot *Protocol) sendAll(data []byte) error {
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
