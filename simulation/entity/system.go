package entity

import (
	"time"
)

type System interface {
	Update(dt time.Duration, e *Entity) error
}
