package prompts

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func SelectProfile(profiles []data.Profile) *selection.Selection {
	var profileNames []string
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name())
	}

	prompt := selection.New("What profile should this Group have?", selection.Choices(profileNames))
	prompt.PageSize = 5

	return prompt
}

func GetGroupName() *textinput.TextInput {
	prompt := textinput.New("What should this Group be named?")
	prompt.Placeholder = "Group Name cannot be empty"
	return prompt
}

func ShouldBeCaptain() *confirmation.Confirmation {
	return confirmation.New("Should this Group be a captain?", confirmation.No)
}

func NewGroup(availableProfiles []data.Profile, availableTraits []data.Trait) (data.Group, error) {
	name, err := GetGroupName().RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	profileChoice, err := SelectProfile(availableProfiles).RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	makeCaptain, err := ShouldBeCaptain().RunPrompt()
	if err != nil {
		return data.Group{}, err
	}

	group, err := data.NewGroup(name, profileChoice.String, availableProfiles)
	if err != nil {
		return data.Group{}, err
	}

	if makeCaptain {
		err = group.MakeCaptain("", availableTraits)
		if err != nil {
			return data.Group{}, err
		}
	}

	return group, nil
}
