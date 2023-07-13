package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fabiandes/slb/demo/simulation/entity"
	zmq "github.com/pebbe/zmq4"
)

type Client struct {
	s *Scene
}

func New(s *Scene) *Client {
	c := &Client{
		s: s,
	}

	return c
}

func (c *Client) Listen(ctx context.Context, addr string) error {
	zctx, err := zmq.NewContext()
	if err != nil {
		return fmt.Errorf("failed to create zmq context: %v", err)
	}

	s, err := zctx.NewSocket(zmq.SUB)
	if err != nil {
		return fmt.Errorf("failed to create socket: %v", err)
	}

	if err := s.Connect(fmt.Sprintf("tcp://%s", addr)); err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	fmt.Println("Client successfully connected!")
	if err := s.SetSubscribe("Entity Update"); err != nil {
		return fmt.Errorf("failed to establish message filer: %v", err)
	}
	fmt.Println("Waiting for updates.")

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancellation")
		default:
			_, err := s.Recv(0)
			if err != nil {
				return fmt.Errorf("failed to receive topic declaration: %v", err)
			}

			msg, err := s.Recv(0)
			if err != nil {
				return fmt.Errorf("failed to receive entity update: %v", err)
			}

			es := []*entity.Entity{}
			if err := json.Unmarshal([]byte(msg), &es); err != nil {
				return fmt.Errorf("failed to unmarshall entities: %v", err)
			}

			fmt.Println("Got update!")

			c.s.Update(es)
		}

	}
}
