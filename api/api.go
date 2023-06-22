package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/fabiandes/spatial-load-balancer/simulation"
	"go.uber.org/zap"
)

type API struct {
	clients []*Client
	mu      sync.Mutex
	logger  *zap.SugaredLogger
}

type APIOptions struct {
	Logger *zap.SugaredLogger
}

// New creates an API with the provided options.
func New(opts *APIOptions) *API {
	api := &API{
		clients: []*Client{},
		mu:      sync.Mutex{},
		logger:  opts.Logger,
	}

	return api
}

// Listen creates a TCP listener at the proved addresses and accepts new clients.
func (api *API) Listen(ctx context.Context, addr string) error {
	// Create a TCP listener and bind it to the provided address.
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen at the proved address: %v", err)
	}

	api.logger.Infof("API is listening at %q", l.Addr())

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("the provided context was cancelled: %v", ctx.Err())
		default:
			conn, err := l.Accept()
			if err != nil {
				return fmt.Errorf("failed to accept connection: %v", err)
			}

			api.logger.Infof("A new connection has been established @ %q", conn.RemoteAddr())

			// Create a client from the connection.
			c := NewClient(conn)
			go api.Handle(ctx, c) // ! We might need to rethink this context.
		}
	}
}

// Handle manages a clients connections.
func (api *API) Handle(ctx context.Context, c *Client) {
	// Register the client
	api.mu.Lock()
	api.clients = append(api.clients, c)
	api.mu.Unlock()

	api.logger.Infow("A new client has been registered", "id", c.id, "addr", c.conn.RemoteAddr())

	defer func() {
		api.mu.Lock()
		defer api.mu.Unlock()

		for i := 0; i < len(api.clients); i++ {
			client := api.clients[i]
			if client.id == c.id {
				api.clients = append(api.clients[:i], api.clients[i+1:]...)
				return
			}
		}
		api.logger.Fatalw("tried to remove client but it did not exist", "client id", c.id)
	}()

	// Continuously read from the clients connection to ensure it is not closed.
	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := make([]byte, 1024)
			if _, err := c.conn.Read(b); err != nil {
				if err == io.EOF {
					api.logger.Warnw("A client has disconnected from the api", "addr", c.conn.RemoteAddr())
					return
				}

				api.logger.Fatalf("unexpected error when reading from client: %v", err)
			}
		}
	}
}

func (api *API) Broadcast(es []*simulation.Entity) error {
	for _, c := range api.clients {
		api.logger.Infof("Sending message to client %q", c.id)
		if err := c.Send(es); err != nil {
			return fmt.Errorf("failed to send entity update to client %q: %v", c.id, err)
		}
		api.logger.Infof("Successfully sent message to client %q", c.id)
	}

	return nil
}
