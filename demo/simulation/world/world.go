package world

import (
	"fmt"

	"github.com/fabiandes/slb/demo/simulation/entity"
)

type World struct {
	Entities    []*entity.Entity
	Width       int
	Height      int
	subscribers []chan []*entity.Entity
}

func NewWorld(es []*entity.Entity, width int, height int) *World {
	w := &World{
		Entities: es,
		Width:    width,
		Height:   height,
	}

	return w
}

// Broadcast send the current world state to all subscribers.
func (w *World) Broadcast() {
	for _, ch := range w.subscribers {
		ch <- w.Entities
	}
}

// Subscribe provides a channel which broadcasts any updates made to the world.
func (w *World) Subscribe() chan []*entity.Entity {
	ch := make(chan []*entity.Entity)
	w.subscribers = append(w.subscribers, ch)
	fmt.Printf("World currently has %d subscribers.\n", len(w.subscribers))

	return ch
}

// Unsubscribe removes the channel from the worlds subscribers.
func (w *World) Unsubscribe(ch chan []*entity.Entity) {
	for i := 0; i < len(w.subscribers); i++ {
		if w.subscribers[i] == ch {
			w.subscribers = append(w.subscribers[:i], w.subscribers[i+1:]...)
			fmt.Printf("World currently has %d subscribers.\n", len(w.subscribers))
			return
		}
	}
}
