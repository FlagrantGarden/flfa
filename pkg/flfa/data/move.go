package data

import (
	"fmt"
)

type Move struct {
	Activation int
	Distance   int
}

func (move *Move) String() string {
	return fmt.Sprintf("%d+ / %d\"", move.Activation, move.Distance)
}
