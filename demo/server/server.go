package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fabiandes/slb/demo/simulation/world"
	zmq "github.com/pebbe/zmq4"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
)

// Name for OpenTelemetry Tracer.
const name = "demo/server"

// Server provides methods for accepting TCP connections and broadcasting world updates to clients.
type Server struct {
	w *world.World
	l *otelzap.SugaredLogger
}

// New creates a new Server.
func New(l *otelzap.SugaredLogger, w *world.World) *Server {
	s := &Server{
		l: l,
		w: w,
	}

	return s
}

// Listen starts the server and begins broadcasting world updates to clients.
func (s *Server) Listen(ctx context.Context) error {
	s.l.Infoln("Starting server...")
	zctx, err := zmq.NewContext()
	if err != nil {
		return fmt.Errorf("failed to create zmq context: %v", err)
	}

	sock, err := zctx.NewSocket(zmq.PUB)
	if err != nil {
		return fmt.Errorf("failed to create socket: %v", err)
	}

	if err := sock.Bind("tcp://*:8080"); err != nil {
		return fmt.Errorf("failed to bind socket: %v", err)
	}

	s.l.Infoln("Server started successfully!")

	if err := s.Broadcast(ctx, sock); err != nil {
		return fmt.Errorf("failed to broadcast: %v", err)
	}

	return nil
}

// Broadcast sends world updates to all connected clients.
func (s *Server) Broadcast(ctx context.Context, sock *zmq.Socket) error {
	s.l.Infoln("Starting broadcast.")
	for es := range s.w.Subscribe() {
		ctx, span := otel.Tracer(name).Start(ctx, "Broadcast")

		// Send the message to all connected clients.
		b, err := json.Marshal(es)
		if err != nil {
			return fmt.Errorf("failed to marshall entities: %v", err)
		}

		s.l.Ctx(ctx).Infow("Successfully marshalled entities", "entities", string(b))

		if _, err := sock.Send("Entity Update", zmq.SNDMORE); err != nil {
			return fmt.Errorf("failed to send topic declaration: %v", err)
		}

		s.l.Ctx(ctx).Infof("Successfully sent topic declaration.")

		if _, err := sock.Send(string(b), 0); err != nil {
			return fmt.Errorf("failed to send entity update: %v", err)
		}

		s.l.Ctx(ctx).Infof("Successfully sent entity update")

		span.End()
	}

	return nil
}
