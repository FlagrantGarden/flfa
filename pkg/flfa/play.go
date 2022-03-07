package flfa

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func SelectProfilePrompt(profiles []data.Profile) *selection.Selection {
	var profileNames []string
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name())
	}

	prompt := selection.New("What profile should this Group have?",
		selection.Choices(profileNames))
	prompt.PageSize = 5

	return prompt
}

func GetGroupNamePrompt() *textinput.TextInput {
	prompt := textinput.New("What should this Group be named?")
	prompt.Placeholder = "Group Name cannot be empty"
	return prompt
}

func ShouldBeCaptainPrompt() *confirmation.Confirmation {
	return confirmation.New("Should this Group be a captain?", confirmation.No)
}

func (ffapi *Api) NewGroupPrompt() (data.Group, error) {
	newGroupNamePrompt := GetGroupNamePrompt()
	name, err := newGroupNamePrompt.RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	newGroupProfilePrompt := SelectProfilePrompt(ffapi.CachedProfiles)
	profileChoice, err := newGroupProfilePrompt.RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	newGroupCaptaincyPrompt := ShouldBeCaptainPrompt()
	makeCaptain, err := newGroupCaptaincyPrompt.RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	group, err := data.NewGroup(name, profileChoice.String, ffapi.CachedProfiles)
	if err != nil {
		return data.Group{}, err
	}

	if makeCaptain {
		err = group.MakeCaptain("", ffapi.CachedTraits)
		if err != nil {
			return data.Group{}, err
		}
	}

	return group, nil
}
