package data

import (
	"fmt"
	"strings"
)

type Missile struct {
	Activation int
	ToHit      int
	Range      int
}

func (missile *Missile) String() string {
	output := strings.Builder{}
	if missile.Activation == 0 {
		output.WriteString("- / - / -")
	} else {
		output.WriteString(fmt.Sprintf("%d+ / %d+ / %d\"", missile.Activation, missile.ToHit, missile.Range))
	}
	return output.String()
}
