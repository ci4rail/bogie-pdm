// Package daprpubsub is a wrapper around the DAPR client
package daprpubsub

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
)

// Connection is the DAPR connection
type Connection struct {
	client      dapr.Client
	nodeID      string
	networkName string
}

// New creates a new instance of DaprPubSub
func New(address string, nodeID string, networkName string) (*Connection, error) {
	client, err := dapr.NewClientWithAddress(address)
	if err != nil {
		return nil, err
	}

	return &Connection{
		client:      client,
		nodeID:      nodeID,
		networkName: networkName,
	}, nil
}

// PubExport publishes a message to the export stream
func (c *Connection) PubExport(subject string, data []byte) error {
	ctx := context.Background()
	subject = fmt.Sprintf("%s.EXPORT.%s", c.nodeID, subject)

	err := c.client.PublishEvent(ctx, c.networkName, subject, data, dapr.PublishEventWithContentType("application/protobuf"))
	return err
}
