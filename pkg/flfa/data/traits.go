package data

type Trait struct {
	Name   string
	Type   string
	Source string
	Roll   int
	Effect string
	Points int
}

func (trait Trait) WithSource(source string) Trait {
	trait.Source = source
	return trait
}

func (trait Trait) WithSubtype(subtype string) Trait {
	trait.Type = subtype
	return trait
}

func FilterTraitsBySource(sourceName string, traitList []Trait) []Trait {
	filteredTraits := []Trait{}
	for _, trait := range traitList {
		if trait.Source == sourceName {
			filteredTraits = append(filteredTraits, trait)
		}
	}
	return filteredTraits
}

func FilterTraitsByType(typeName string, traitList []Trait) []Trait {
	filteredTraits := []Trait{}
	for _, trait := range traitList {
		if trait.Type == typeName {
			filteredTraits = append(filteredTraits, trait)
		}
	}
	return filteredTraits
}

func GetTraitByName(name string, traitList []Trait) Trait {
	for _, trait := range traitList {
		if trait.Name == name {
			return trait
		}
	}
	return Trait{}
}
