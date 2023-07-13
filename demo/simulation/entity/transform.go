package entity

import "github.com/fabiandes/slb/demo/simulation/vector"

type TransformComponent struct {
	Position vector.Vector2
	Rotation vector.Vector2
	Scale    vector.Vector2
}

func NewTransform() *TransformComponent {
	t := &TransformComponent{
		Scale: vector.Vector2{
			X: 1,
			Y: 1,
		},
	}

	return t
}
