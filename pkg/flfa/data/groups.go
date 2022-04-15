package data

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Group struct {
	Name             string
	Id               string
	ProfileName      string
	Melee            Melee
	Move             Move
	Missile          Missile
	FightingStrength FightingStrength `mapstructure:"fighting_strength"`
	Resolve          int
	Toughness        int
	Traits           []string
	Points           int
	Captain          Trait
	Addenda          map[string]any
	// Tags             []Tag
}

type Grouper interface {
	Initialize(availableProfiles []Profile) error
}

func (group *Group) Initialize(availableProfiles []Profile) error {
	profile, err := GetProfile(group.ProfileName, availableProfiles)
	if err != nil {
		return err
	}

	group.Melee = profile.Melee
	group.Move = profile.Move
	group.Missile = profile.Missile
	group.FightingStrength = profile.FightingStrength
	group.Resolve = profile.Resolve
	group.Toughness = profile.Toughness
	group.Traits = profile.Traits
	group.Points = profile.Points
	group.Addenda = make(map[string]any)

	return nil
}

func (group Group) ToSlice() (groups []Group) {
	groups = append(groups, group)
	return groups
}

func (group *Group) DemoteFromCaptain() {
	group.Captain = Trait{}
}

func (group *Group) PromoteToCaptain(trait *Trait, availableTraits ...Trait) error {
	if trait != nil {
		group.Captain = *trait
		return nil
	}

	rolledTrait, err := RollForCaptainTrait(availableTraits)
	if err != nil {
		return fmt.Errorf("unable to make %s into a captain: %s", group.Name, err)
	}

	log.Trace().Msgf("Rolled a %d for Captain's trait for %s", rolledTrait.Roll, group.Name)
	group.Captain = rolledTrait

	return nil
}

func NewGroup(name string, profileName string, availableProfiles []Profile) (Group, error) {
	group := Group{
		Name:        name,
		ProfileName: profileName,
	}

	err := group.Initialize(availableProfiles)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}
