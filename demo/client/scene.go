package client

import (
	"gioui.org/app"
	"github.com/fabiandes/slb/demo/simulation/entity"
)

type Scene struct {
	Entities []*entity.Entity
	Window   *app.Window
}

func (s *Scene) Update(es []*entity.Entity) {
	s.Entities = es
	s.Window.Invalidate()
}
