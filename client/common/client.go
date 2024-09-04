package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   *Protocol
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}

	setupSignalHandler(client)

	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := CreateConnection(c.config.ServerAddress)

	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		clientMessage := CreateClientMessage(
			c.config.ID,
			os.Getenv("NOMBRE"),
			os.Getenv("APELLIDO"),
			os.Getenv("DOCUMENTO"),
			os.Getenv("NACIMIENTO"),
			os.Getenv("NUMERO"),
		)

		err := c.conn.Send(*clientMessage)

		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		msg, err := c.conn.ReceiveAndCloseConnection()

		if err != nil {
			log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		} else {
			log.Infof("action: receive_response | result: success | message: %v", msg.msg)
		}

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			clientMessage.document,
			clientMessage.number,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func handleSignals(sigs chan os.Signal, client *Client) {
	sig := <-sigs
	log.Infof("action: client_closed | result: in_progress | client_id: %v | signal: %v", client.config.ID, sig)
	if client.conn != nil {
		client.conn.CloseConnection()
	}
	log.Infof("action: client_closed | result: success")
	os.Exit(0)
}

func setupSignalHandler(client *Client) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go handleSignals(sigs, client)
}
