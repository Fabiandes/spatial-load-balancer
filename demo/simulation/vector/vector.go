package vector

import (
	"log"
	"math"
)

type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v *Vector2) Normalize() Vector2 {
	m := v.Magnitude()
	if m == 0 {
		log.Fatal("Cannot normalize the 0 vector.")
	}

	return Vector2{
		X: v.X / m,
		Y: v.Y / m,
	}
}
