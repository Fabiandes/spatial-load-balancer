package api

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/fabiandes/spatial-load-balancer/simulation"
	"github.com/google/uuid"
)

type Client struct {
	id   string
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		conn: conn,
	}

	// Generate a random ID for the client.
	c.id = uuid.NewString()

	return c
}

func (c *Client) Send(es []*simulation.Entity) error {
	// Marshall the message
	encoder := gob.NewEncoder(c.conn)
	if err := encoder.Encode(es); err != nil {
		return fmt.Errorf("failed to encode entities: %v", err)
	}

	return nil
}
