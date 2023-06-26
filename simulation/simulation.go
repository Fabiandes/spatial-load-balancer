package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/fabiandes/spatial-load-balancer/simulation/component"
	"github.com/fabiandes/spatial-load-balancer/simulation/entity"
	"github.com/fabiandes/spatial-load-balancer/simulation/system/worker"
	"github.com/fabiandes/spatial-load-balancer/simulation/vector"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	WorldWidth  = 1000
	WorldHeight = 1000
)

type Simulation struct {
	logger         *zap.SugaredLogger
	entities       []*entity.Entity
	subscribers    []chan []*entity.Entity
	simulationRate int
}

type Options struct {
	StartingEntityCount int
	SimulationRate      int
	Logger              *zap.SugaredLogger
}

func New(opts *Options) (*Simulation, error) {
	s := &Simulation{
		logger:         opts.Logger,
		simulationRate: opts.SimulationRate,
	}

	// Generate a set of Entities within the world.
	es := make([]*entity.Entity, opts.StartingEntityCount)
	s.logger.Infof("Generating %d entities...", opts.StartingEntityCount)
	for i := 0; i < opts.StartingEntityCount; i++ {
		// Generate unique ID for the Entity.
		id := uuid.NewString()

		// Create Transform component and move Entity to a random position within the world.
		t := component.NewTransform()
		t.Position = vector.Vector{
			X: rand.Float64() * WorldWidth,
			Y: rand.Float64() * WorldHeight,
		}

		// Create Locomotion component.
		l := component.Locomotion{
			Speed: 1.5,
		}

		// Create an Entity and attach Systems to it.
		e := entity.NewEntity(id, t, l, s.logger)

		w := worker.NewWorker(s.logger)
		e.Attach(w)

		es[i] = e
		s.logger.Infow("Successfully generated an entity", "entity", e)
	}
	s.entities = es
	s.logger.Infoln("Successfully generated all entities")

	return s, nil
}

func (s *Simulation) Update(ctx context.Context) error {
	start := time.Now()
	// Update each entity by calling all systems.
	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < len(s.entities); i++ {
		e := s.entities[i]
		g.Go(e.Update)
	}

	// Wait for all updates to complete.
	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to update all entities: %v", err)
	}

	// ? We could use a go routine for this but we could need to use mutex locking.
	// Publish changes to all subscribers
	s.Publish(ctx)

	s.logger.Infow("Update completed", "duration", time.Since(start))
	return nil
}

func (s *Simulation) Publish(ctx context.Context) {
	for _, ch := range s.subscribers {
		ch <- s.entities
	}
}

func (s *Simulation) Subscribe(ch chan []*entity.Entity) {
	s.subscribers = append(s.subscribers, ch)
}

func (s *Simulation) Unsubscribe(ch chan []*entity.Entity) {
	for i := 0; i < len(s.subscribers); i++ {
		if s.subscribers[i] == ch {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			return
		}
	}
}

func (s *Simulation) Run(ctx context.Context) error {
	t := time.NewTicker(time.Second / time.Duration(s.simulationRate))

	for {
		select {
		case <-t.C:
			s.Update(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}
