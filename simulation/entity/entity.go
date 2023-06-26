package entity

import (
	"time"

	"github.com/fabiandes/spatial-load-balancer/simulation/component"
	"go.uber.org/zap"
)

type Entity struct {
	Id         string
	Transform  component.Transform
	Locomotion component.Locomotion
	logger     *zap.SugaredLogger
	systems    []System
}

func NewEntity(id string, t component.Transform, l component.Locomotion, logger *zap.SugaredLogger) *Entity {
	e := &Entity{
		Id:         id,
		Transform:  t,
		Locomotion: l,
		logger:     logger,
		systems:    []System{},
	}

	return e
}

func (e *Entity) Update() error {
	for i := 0; i < len(e.systems); i++ {
		s := e.systems[i]
		// TODO: We could also spawn go routines to handle this but we should investigate performance first.
		s.Update(time.Second, e)
	}
	return nil
}

func (e *Entity) Attach(s System) {
	e.logger.Infow("Attaching new system to entity.", "entity id", e.Id, "system", s)
	e.systems = append(e.systems, s)
}
