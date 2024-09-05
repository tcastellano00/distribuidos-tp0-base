package common

import (
	"fmt"
	"strings"
)

type ClientMessage struct {
	agency    string
	firstName string
	lastName  string
	document  string
	birthdate string
	number    string
}

type ServerMessage struct {
	success bool
	msg     string
}

func CreateClientMessage(
	agency string, firstName string, lastName string, document string,
	birthdate string, number string) *ClientMessage {

	message := &ClientMessage{
		agency:    agency,
		firstName: firstName,
		lastName:  lastName,
		document:  document,
		birthdate: birthdate,
		number:    number,
	}

	return message
}

func CreateServerMessage(
	success bool, msg string) *ServerMessage {

	message := &ServerMessage{
		success: success,
		msg:     msg,
	}

	return message
}

func (msg *ClientMessage) Serialize() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s\n",
		msg.agency,
		msg.firstName,
		msg.lastName,
		msg.document,
		msg.birthdate,
		msg.number,
	)
}

func DeserializeServerMessage(serverMessage string) *ServerMessage {
	messageParts := strings.Split(strings.TrimSpace(serverMessage), ",")

	return CreateServerMessage(
		messageParts[0] == "OK",
		messageParts[1],
	)
}
