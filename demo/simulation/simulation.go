package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/fabiandes/slb/demo/simulation/vector"
	"github.com/fabiandes/slb/demo/simulation/world"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
)

const name = "demo/simulation"

const EntitySpeed = 1.54  // The movement speed of an entity.
const SimulationRate = 30 // Number of updates per second
const TickDuration = time.Second / SimulationRate

type Simulation struct {
	l *otelzap.SugaredLogger
	w *world.World
}

func New(l *otelzap.SugaredLogger, w *world.World) *Simulation {
	s := &Simulation{
		l: l,
		w: w,
	}

	return s
}

// Simulate begins the Simulation and starts the simulation loop.
func (s *Simulation) Simulate(ctx context.Context) error {
	s.l.Infoln("Starting simulation.")
	t := time.NewTicker(TickDuration)

	for {
		select {
		case <-t.C:
			if err := s.Update(ctx, TickDuration); err != nil {
				return fmt.Errorf("failed to perform update: %v", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("context cancellation")
		}
	}
}

// Update performs a step in the simulation.
func (s *Simulation) Update(ctx context.Context, dt time.Duration) error {
	defer s.w.Broadcast()

	ctx, span := otel.Tracer(name).Start(ctx, "Update")
	defer span.End()

	start := time.Now()
	for _, e := range s.w.Entities {
		s.l.Ctx(ctx).Infow("Updating entity", "entity id", e.Id)

		// Move towards the destination.
		distRemaining := EntitySpeed * dt.Seconds()
		for distRemaining > 0 {
			s.l.Ctx(ctx).Infow("Moving entity", "entity id", e.Id, "distance remaining", distRemaining)
			// Generate a destination if the entity does not currently have one.
			if e.Navigation.Destination == nil {
				d := &vector.Vector2{
					X: rand.Float64() * float64(s.w.Width),
					Y: rand.Float64() * float64(s.w.Height),
				}
				e.Navigation.Destination = d

				s.l.Ctx(ctx).Infow("Generated destination for entity", "entity id", e.Id, "destination", d)
			}

			// Generate a set of waypoints if the entity does not currently have any.
			if len(e.Navigation.Waypoints) == 0 {
				// TODO: Actually generate a path with path finding.
				w := []*vector.Vector2{e.Navigation.Destination}
				e.Navigation.Waypoints = w
				s.l.Ctx(ctx).Infow("Generated waypoints for entity", "entity id", e.Id, "waypoints", w)
			}

			w := e.Navigation.Waypoints[0]
			displacement := vector.Vector2{
				X: w.X - e.Transform.Position.X,
				Y: w.Y - e.Transform.Position.Y,
			}

			// If the entity has arrived at a waypoint, or can reach it, discard it and continue moving.
			if displacement.Magnitude() == 0 || displacement.Magnitude() < distRemaining {
				e.Transform.Position = *w
				e.Navigation.Waypoints = e.Navigation.Waypoints[1:]
				distRemaining -= displacement.Magnitude()
				s.l.Ctx(ctx).Infow("Entity moved to waypoint", "entity id", e.Id, "distance traveled", displacement.Magnitude())

				if w == e.Navigation.Destination {
					s.l.Ctx(ctx).Infow("Entity reach destination", "entity id", e.Id)
					e.Navigation.Destination = nil
				}
				continue
			}

			// Move as far towards the next waypoint as possible.
			dir := displacement.Normalize()
			move := vector.Vector2{
				X: dir.X * distRemaining,
				Y: dir.Y * distRemaining,
			}
			e.Transform.Position.X += move.X
			e.Transform.Position.Y += move.Y
			distRemaining = 0
			s.l.Ctx(ctx).Infow("Entity moved towards waypoint", "entity id", e.Id)
		}
	}

	if dur := time.Since(start); dur > TickDuration {
		s.l.Ctx(ctx).Warnw("Simulation took too long", "duration", dur, "expected duration", TickDuration)
	}
	return nil
}
