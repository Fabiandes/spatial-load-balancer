package worker

import (
	"math/rand"
	"time"

	"github.com/fabiandes/spatial-load-balancer/simulation/entity"
	"github.com/fabiandes/spatial-load-balancer/simulation/vector"
)

type Job interface {
	PerformNextTask(dt time.Duration, e *entity.Entity) error
}

type Wandering struct {
	waypoints   []vector.Vector
	destination *vector.Vector
	worldWidth  int
	worldHeight int
}

func WanderingJob(worldWidth, worldHeight int) *Wandering {
	w := &Wandering{
		waypoints:   []vector.Vector{},
		worldWidth:  worldWidth,
		worldHeight: worldHeight,
	}

	return w
}

func (j *Wandering) PerformNextTask(dt time.Duration, e *entity.Entity) error {
	// Generate a set of waypoints if the work doesn't have any.
	if len(j.waypoints) == 0 {
		j.waypoints = j.Waypoints()
	}

	return nil
}

func (j *Wandering) Waypoints() []vector.Vector {
	// Generate a destination if the work doesn't have one.
	if j.destination == nil {
		j.destination = j.Destination()
	}

	// TODO: Generate a path to the destination using A* path finding.
	w := []vector.Vector{}

	return w
}

func (j *Wandering) Destination() *vector.Vector {
	d := &vector.Vector{
		X: rand.Float64() * float64(j.worldWidth),
		Y: rand.Float64() * float64(j.worldHeight),
	}

	return d
}
