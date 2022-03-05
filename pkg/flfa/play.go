package flfa

import (
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func (ffapi *Api) SelectProfilePrompt() *selection.Selection {
	var validProfiles []string
	for _, profile := range ffapi.CachedProfiles {
		validProfiles = append(validProfiles, profile.Name())
	}

	prompt := selection.New("What profile should this Group have?",
		selection.Choices(validProfiles))
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

func (ffapi *Api) NewGroupPrompt() (Group, error) {
	newGroupNamePrompt := GetGroupNamePrompt()
	name, err := newGroupNamePrompt.RunPrompt()
	if err != nil {
		return Group{}, err
	}

	newGroupProfilePrompt := ffapi.SelectProfilePrompt()
	profileChoice, err := newGroupProfilePrompt.RunPrompt()
	if err != nil {
		return Group{}, err
	}

	newGroupCaptaincyPrompt := ShouldBeCaptainPrompt()
	makeCaptain, err := newGroupCaptaincyPrompt.RunPrompt()
	if err != nil {
		return Group{}, err
	}

	group, err := ffapi.NewGroup(name, profileChoice.String)
	if err != nil {
		return Group{}, err
	}

	if makeCaptain {
		err = group.MakeCaptain("")
		if err != nil {
			return Group{}, err
		}
	}

	return group, nil
}
