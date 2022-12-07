package nats

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
)

// Connection is the NATS connection
type Connection struct {
	nc     *nats.Conn
	js     nats.JetStream
	nodeID string
}

// Connect connects to the NATS server and jetstream
func Connect(address string, nodeID string, creds string) (*Connection, error) {
	options := []nats.Option{}
	if creds != "" {
		options = append(options, nats.UserCredentials(creds))
	}
	nc, err := nats.Connect(address, options...)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Connection{
		nc:     nc,
		js:     js,
		nodeID: nodeID,
	}, nil
}

// PubExport publishes a message to the jetstream subject
func (c *Connection) PubExport(subject string, data []byte) error {
	subject = fmt.Sprintf("%s.EXPORT.%s", c.nodeID, subject)
	_, err := c.js.Publish(subject, data)
	return err
}
