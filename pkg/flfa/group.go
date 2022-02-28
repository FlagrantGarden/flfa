package flfa

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/prompt"
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

func (ffapi *Api) NewGroupPrompt() (Group, error) {
	var validProfiles []string
	for _, profile := range ffapi.CachedProfiles {
		validProfiles = append(validProfiles, profile.Name())
	}
	profilePrompt := prompt.PromptContent{
		ErrorMessage: fmt.Sprintf("Please choose a valid profile from this list: %s", validProfiles),
		Label:        "What base profile should this Group have?",
	}
	err := profilePrompt.GetSelection(validProfiles)
	if err != nil {
		return Group{}, err
	}

	namePrompt := prompt.PromptContent{
		ErrorMessage: "Please provide a name.",
		Label:        "What name should this Group be called?",
	}
	err = namePrompt.GetInput()
	if err != nil {
		return Group{}, err
	}

	group, err := ffapi.NewGroup(namePrompt.Result, profilePrompt.Result)
	if err != nil {
		return Group{}, err
	}

	captainPrompt := prompt.PromptContent{
		Label: "Should this Group be a Captain?",
	}
	captainPrompt.GetConfirmation()

	if captainPrompt.Result == "yes" {
		err := group.MakeCaptain("")
		if err != nil {
			return Group{}, err
		}
	}

	return group, nil
}
