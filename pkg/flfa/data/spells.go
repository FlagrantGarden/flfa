package data

type Spell struct {
	Source   string
	Name     string
	Check    int
	Range    int
	Target   string
	Duration string
	Effect   string
}

func (spell Spell) WithSource(source string) Spell {
	spell.Source = source
	return spell
}
