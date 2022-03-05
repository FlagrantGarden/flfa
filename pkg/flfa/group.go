package flfa

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
	BaseProfileName  string
	Melee            Melee
	Move             Move
	Missile          Missile
	FightingStrength FightingStrength
	Resolve          int
	Toughness        int
	Traits           []string
	Points           int
	Captain          Trait
	Api              *Api
	// Tags             []Tag
}

type Grouper interface {
	InitializeToBaseProfile() error
}

func (g *Group) InitializeToBaseProfile() error {
	baseProfile, err := GetBaseProfile(g.BaseProfileName, g.Api.CachedProfiles)
	if err != nil {
		return err
	}

	g.Melee = baseProfile.Melee
	g.Move = baseProfile.Move
	g.Missile = baseProfile.Missile
	g.FightingStrength = baseProfile.FightingStrength
	g.Resolve = baseProfile.Resolve
	g.Toughness = baseProfile.Toughness
	g.Traits = baseProfile.Traits
	g.Points = baseProfile.Points

	return nil
}

func (g *Group) MakeCaptain(traitName string) error {
	rand.Seed(time.Now().UnixNano())
	result, _, err := dice.Roll("3d6")
	if err != nil {
		return fmt.Errorf("unable to make %s into a captain: %s", g.Name, err)
	}
	log.Trace().Msgf("Rolled a %d for Captain's trait for %s", result.Int(), g.Name)
	for _, trait := range FilterTraitsByType("Captain", g.Api.CachedTraits) {
		if traitName != "" && traitName == trait.Name {
			g.Captain = trait
			return nil
		} else if trait.Roll == result.Int() {
			g.Captain = trait
			return nil
		}
	}
	return fmt.Errorf("unable to make %s into a captain: could not assign trait '%s'", g.Name, traitName)
}

func (g *Group) MarkdownTableEntry() string {
	output := strings.Builder{}
	traits := g.Traits
	if g.Captain.Name == "" {
		output.WriteString(fmt.Sprintf("| %s |", g.Name))
	} else {
		output.WriteString(fmt.Sprintf("| **%s** |", g.Name))
		traits = append([]string{fmt.Sprintf("**%s**", g.Captain.Name)}, traits...)
	}
	output.WriteString(fmt.Sprintf(" %s |", g.BaseProfileName))
	output.WriteString(fmt.Sprintf(" %s |", g.Melee.String()))
	output.WriteString(fmt.Sprintf(" %s |", g.Missile.String()))
	output.WriteString(fmt.Sprintf(" %s |", g.Move.String()))
	output.WriteString(fmt.Sprintf(" %s |", g.FightingStrength.String()))
	output.WriteString(fmt.Sprintf(" %d+ |", g.Resolve))
	output.WriteString(fmt.Sprintf(" %d |", g.Toughness))
	output.WriteString(fmt.Sprintf(" %s |\n", strings.Join(traits, ", ")))
	return output.String()
}

func (g *Group) JSON() string {
	data := map[string]interface{}{
		"name":            g.Name,
		"id":              g.Id,
		"baseProfileName": g.BaseProfileName,
		"melee": map[string]interface{}{
			"activation":     g.Melee.Activation,
			"toHitAttacking": g.Melee.ToHitAttacking,
			"toHitDefending": g.Melee.ToHitDefending,
		},
		"move": map[string]interface{}{
			"activation": g.Move.Activation,
			"distance":   g.Move.Distance,
		},
		"fightingStrength": map[string]interface{}{
			"current": g.FightingStrength.Current,
			"maximum": g.FightingStrength.Maximum,
		},
		"resolve":   g.Resolve,
		"toughness": g.Toughness,
		"traits":    g.Traits,
	}
	if g.Missile.Activation == 0 {
		data["missile"] = map[string]interface{}{}
	} else {
		data["missile"] = map[string]interface{}{
			"activation": g.Missile.Activation,
			"toHit":      g.Missile.ToHit,
			"range":      g.Missile.Range,
		}
	}
	if g.Captain.Name == "" {
		data["captain"] = map[string]interface{}{}
	} else {
		data["captain"] = map[string]interface{}{
			"name": g.Captain.Name,
		}
	}
	jsonString, _ := json.Marshal(data)
	return string(jsonString)
}

func (ffapi *Api) NewGroup(name string, profileName string) (Group, error) {
	group := Group{
		Name:            name,
		BaseProfileName: profileName,
		Api:             ffapi,
	}

	err := group.InitializeToBaseProfile()
	if err != nil {
		return Group{}, err
	}

	return group, nil
}
