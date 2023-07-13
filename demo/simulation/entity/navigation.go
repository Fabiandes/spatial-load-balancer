package entity

import "github.com/fabiandes/slb/demo/simulation/vector"

type NavigationComponent struct {
	Destination *vector.Vector2
	Waypoints   []*vector.Vector2
}
