package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func New(x, y float64) Vector {
	return Vector{
		X: x,
		Y: y,
	}
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v *Vector) Normalize() (Vector, error) {
	if v.Magnitude() == 0 {
		return Vector{}, fmt.Errorf("cannot normalize zero vector")
	}

	m := v.Magnitude()

	return Vector{
		X: v.X / m,
		Y: v.Y / m,
	}, nil
}
