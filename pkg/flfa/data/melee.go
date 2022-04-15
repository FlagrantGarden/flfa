package data

import (
	"fmt"
	"strings"
)

type Melee struct {
	Activation     int
	ToHitAttacking int `mapstructure:"to_hit_attacking"`
	ToHitDefending int `mapstructure:"to_hit_defending"`
}

func (melee *Melee) String() string {
	output := strings.Builder{}
	if melee.Activation == 0 {
		output.WriteString("- / ")
	} else {
		output.WriteString(fmt.Sprintf("%d+ / ", melee.Activation))
	}
	if melee.ToHitAttacking == 0 {
		output.WriteString("- / ")
	} else if melee.ToHitAttacking < 6 {
		output.WriteString(fmt.Sprintf("%d+ / ", melee.ToHitAttacking))
	} else {
		output.WriteString("6 / ")
	}
	if melee.ToHitDefending < 6 {
		output.WriteString(fmt.Sprintf("%d+", melee.ToHitDefending))
	} else {
		output.WriteString("6")
	}
	return output.String()
}
