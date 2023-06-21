package simulation

import (
	"github.com/fabiandes/spatial-load-balancer/simulation/component"
	"github.com/fabiandes/spatial-load-balancer/simulation/system"
	"go.uber.org/zap"
)

type Entity struct {
	Id        string
	Transform component.Transform
	logger    *zap.SugaredLogger
	Systems   []system.System
}

func NewEntity(id string, t component.Transform, logger *zap.SugaredLogger) *Entity {
	e := &Entity{
		Id:        id,
		Transform: t,
		logger:    logger,
		Systems:   []system.System{},
	}

	return e
}

func (e *Entity) Update() error {
	for i := 0; i < len(e.Systems); i++ {
		s := e.Systems[i]
		// TODO: We could also spawn go routines to handle this but we should investigate performance first.
		s.Update()
	}
	return nil
}

func (e *Entity) Attach(s system.System) {
	e.logger.Infow("Attaching new system to entity.", "entity id", e.Id, "system", s)
	e.Systems = append(e.Systems, s)
}
