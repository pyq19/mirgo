package common

import (
	"fmt"
)

type Point struct {
	X uint32
	Y uint32
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}
