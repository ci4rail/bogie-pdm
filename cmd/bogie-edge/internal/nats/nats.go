package nats

import (
	nats "github.com/nats-io/nats.go"
)

// Connection is the NATS connection
type Connection struct {
	nc *nats.Conn
	js nats.JetStream
}

// Connect connects to the NATS server and jetstream
func Connect(address string) (*Connection, error) {
	nc, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Connection{
		nc: nc,
		js: js,
	}, nil
}

// PubJs publishes a message to the jetstream subject
func (c *Connection) PubJs(subject string, data []byte) error {
	_, err := c.js.Publish(subject, data)
	return err
}
