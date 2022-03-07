package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/justinian/dice"
	"github.com/rs/zerolog/log"
)

type Group struct {
	Name             string
	Id               string
	ProfileName      string
	Melee            Melee
	Move             Move
	Missile          Missile
	FightingStrength FightingStrength
	Resolve          int
	Toughness        int
	Traits           []string
	Points           int
	Captain          Trait
	// Tags             []Tag
}

type Grouper interface {
	Initialize(availableProfiles []Profile) error
}

func (group *Group) MarkdownTableEntry() string {
	output := strings.Builder{}
	traits := group.Traits
	if group.Captain.Name == "" {
		output.WriteString(fmt.Sprintf("| %s |", group.Name))
	} else {
		output.WriteString(fmt.Sprintf("| **%s** |", group.Name))
		traits = append([]string{fmt.Sprintf("**%s**", group.Captain.Name)}, traits...)
	}
	output.WriteString(fmt.Sprintf(" %s |", group.ProfileName))
	output.WriteString(fmt.Sprintf(" %s |", group.Melee.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.Missile.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.Move.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.FightingStrength.String()))
	output.WriteString(fmt.Sprintf(" %d+ |", group.Resolve))
	output.WriteString(fmt.Sprintf(" %d |", group.Toughness))
	output.WriteString(fmt.Sprintf(" %s |\n", strings.Join(traits, ", ")))
	return output.String()
}

func (group *Group) JSON() string {
	data := map[string]interface{}{
		"name":        group.Name,
		"id":          group.Id,
		"profileName": group.ProfileName,
		"melee": map[string]interface{}{
			"activation":     group.Melee.Activation,
			"toHitAttacking": group.Melee.ToHitAttacking,
			"toHitDefending": group.Melee.ToHitDefending,
		},
		"move": map[string]interface{}{
			"activation": group.Move.Activation,
			"distance":   group.Move.Distance,
		},
		"fightingStrength": map[string]interface{}{
			"current": group.FightingStrength.Current,
			"maximum": group.FightingStrength.Maximum,
		},
		"resolve":   group.Resolve,
		"toughness": group.Toughness,
		"traits":    group.Traits,
	}
	if group.Missile.Activation == 0 {
		data["missile"] = map[string]interface{}{}
	} else {
		data["missile"] = map[string]interface{}{
			"activation": group.Missile.Activation,
			"toHit":      group.Missile.ToHit,
			"range":      group.Missile.Range,
		}
	}
	if group.Captain.Name == "" {
		data["captain"] = map[string]interface{}{}
	} else {
		data["captain"] = map[string]interface{}{
			"name": group.Captain.Name,
		}
	}
	jsonString, _ := json.Marshal(data)
	return string(jsonString)
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

	return nil
}

func (group *Group) MakeCaptain(traitName string, availableTraits []Trait) error {
	rand.Seed(time.Now().UnixNano())
	result, _, err := dice.Roll("3d6")
	if err != nil {
		return fmt.Errorf("unable to make %s into a captain: %s", group.Name, err)
	}
	log.Trace().Msgf("Rolled a %d for Captain's trait for %s", result.Int(), group.Name)
	for _, trait := range FilterTraitsByType("Captain", availableTraits) {
		if traitName != "" && traitName == trait.Name {
			group.Captain = trait
			return nil
		} else if trait.Roll == result.Int() {
			group.Captain = trait
			return nil
		}
	}
	return fmt.Errorf("unable to make %s into a captain: could not assign trait '%s'", group.Name, traitName)
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
