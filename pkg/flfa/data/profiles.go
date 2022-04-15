package data

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Profile struct {
	Source           string
	Type             string
	Category         string
	Melee            Melee
	Move             Move
	Missile          Missile
	FightingStrength FightingStrength `mapstructure:"fighting_strength"`
	Resolve          int
	Toughness        int
	Traits           []string
	Points           int
}

type Profiler interface {
	Name() string
}

func (profile *Profile) Name() string {
	return fmt.Sprintf("%s %s", profile.Type, profile.Category)
}

func (profile Profile) WithSource(source string) Profile {
	profile.Source = source
	return profile
}

func GetProfile(name string, profileList []Profile) (Profile, error) {
	log.Trace().Msgf("searching for profile '%s'", name)
	for _, profile := range profileList {
		if profile.Name() == name {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("no profile found that matches name '%s'", name)
}

func GetProfileTraits(profile Profile, traitList []Trait) []Trait {
	traits := []Trait{}
	for _, traitName := range profile.Traits {
		trait := GetTraitByName(traitName, traitList)
		if trait.Name != "" {
			traits = append(traits, trait)
		}
	}
	return traits
}
