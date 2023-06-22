package renderer

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"gioui.org/app"
	"github.com/fabiandes/spatial-load-balancer/simulation"
)

type Scene struct {
	Entities []*simulation.Entity
	window   *app.Window
}

func NewScene(w *app.Window) *Scene {
	s := &Scene{
		Entities: []*simulation.Entity{},
		window:   w,
	}

	return s
}

func (s *Scene) Initialize() error {
	// Connect to the simulation.
	fmt.Println("Connecting to simulation...")
	conn, err := net.Dial("tcp4", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to simulation")

	// Update the window whenever an update is received.
	for {
		// Unmarshall the received data.
		fmt.Println("Waiting for update...")
		es := new([]*simulation.Entity)
		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(es); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully unmarshalled update")

		// Update the scene and trigger a re-render.
		s.Entities = *es
		s.window.Invalidate()
	}
}
