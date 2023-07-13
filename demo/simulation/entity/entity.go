package entity

type Entity struct {
	Id         string
	Transform  *TransformComponent
	Navigation *NavigationComponent
}

func NewEntity(t *TransformComponent, n *NavigationComponent) *Entity {
	e := &Entity{
		Transform:  t,
		Navigation: n,
	}

	return e
}
