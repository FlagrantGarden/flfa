package data

import "fmt"

type FightingStrength struct {
	Current int
	Maximum int
}

func (fightingStrength *FightingStrength) String() string {
	return fmt.Sprintf("%d / %d", fightingStrength.Current, fightingStrength.Maximum)
}
