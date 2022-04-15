package prompts

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	tprompts "github.com/FlagrantGarden/flfa/pkg/flfa/tui/traits/prompts"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func SelectProfile(profiles []data.Profile) *selection.Selection {
	var profileNames []string
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name())
	}

	return selector.NewStringSelector(
		"What profile should this Group have?",
		profileNames,
		selector.WithPageSize(5),
	)
}

func SelectProfileModel(profiles []data.Profile) *selection.Model {
	return selection.NewModel(SelectProfile(profiles))
}

func GetGroupName() *textinput.TextInput {
	return texter.New(
		"What should this Group be named?",
		texter.WithPlaceholder("Group Name cannot be empty"),
		texter.WithInputWidth(30),
	)
}

func GetGroupNameModel() *textinput.Model {
	return textinput.NewModel(GetGroupName())
}

func SelectGroupEditingOption(hasAddableTraits bool, hasRemovableTraits bool) *selection.Selection {
	options := []string{
		"Save & Continue",
		"Change Name",
		"Change Base Profile",
	}

	if hasAddableTraits {
		options = append(options, "Add Special Trait")
	}

	if hasRemovableTraits {
		options = append(options, "Remove Special Trait")
	}

	return selector.NewStringSelector(
		"What would you like to do with this Group?",
		options,
		selector.WithPageSize(5),
	)
}

func SelectGroupEditingOptionModel(hasAddableTraits bool, hasRemovableTraits bool) *selection.Model {
	return selection.NewModel(SelectGroupEditingOption(hasAddableTraits, hasRemovableTraits))
}

func SelectAddSpecialTrait(applicableTraits []data.Trait) *selection.Selection {
	return tprompts.SelectTrait(
		"Which trait do you want to add?",
		applicableTraits,
	)
}

func SelectAddSpecialTraitModel(applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectAddSpecialTrait(applicableTraits))
}

func SelectRemoveSpecialTrait(applicableTraits []data.Trait) *selection.Selection {
	return tprompts.SelectTrait(
		"Which trait do you want to remove?",
		applicableTraits,
	)
}

func SelectRemoveSpecialTraitModel(applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectRemoveSpecialTrait(applicableTraits))
}

func ConfirmChangeBaseProfile(new_profile string) *confirmation.Confirmation {
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("Are you sure you want to update the base profile to '%s'?\n", new_profile))
	messageBuilder.WriteString("Doing so will reset this group to that base profile. ")
	messageBuilder.WriteString("You'll lose any changes other than the group's name and captaincy.")
	message := lipgloss.NewStyle().Width(120).Render(messageBuilder.String())
	message = strings.Join([]string{message, ""}, "\n")

	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmChangeBaseProfileModel(new_profile string) *confirmation.Model {
	return confirmation.NewModel(ConfirmChangeBaseProfile(new_profile))
}
