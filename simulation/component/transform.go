package component

import (
	"github.com/fabiandes/spatial-load-balancer/simulation/vector"
)

type Transform struct {
	Position vector.Vector
	Rotation vector.Vector
	Scale    vector.Vector
}

func NewTransform() Transform {
	t := Transform{
		Scale: vector.Vector{
			X: 1,
			Y: 1,
		},
	}
	return t
}
