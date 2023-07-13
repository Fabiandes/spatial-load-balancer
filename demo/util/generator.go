package util

import (
	"math/rand"
	"time"

	"github.com/fabiandes/slb/demo/simulation/entity"
	"github.com/fabiandes/slb/demo/simulation/vector"
	"github.com/fabiandes/slb/demo/simulation/world"
)

func GenerateWorld(entityCount int, mapWidth int, mapHeight int) *world.World {
	rand.Seed(time.Now().UnixNano())
	es := make([]*entity.Entity, entityCount)
	for i := 0; i < entityCount; i++ {
		t := entity.NewTransform()
		t.Position = vector.Vector2{
			X: rand.Float64() * float64(mapWidth),
			Y: rand.Float64() * float64(mapHeight),
		}

		n := &entity.NavigationComponent{}

		es[i] = entity.NewEntity(t, n)
	}

	w := world.NewWorld(es, mapWidth, mapHeight)

	return w
}
