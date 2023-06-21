package component

import (
	"github.com/fabiandes/spatial-load-balancer/simulation/util"
)

type Transform struct {
	Position util.Vector
	Rotation util.Vector
	Scale    util.Vector
}

func NewTransform() Transform {
	t := Transform{
		Scale: util.Vector{
			X: 1,
			Y: 1,
		},
	}
	return t
}
