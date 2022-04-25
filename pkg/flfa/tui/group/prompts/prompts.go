package prompts

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	tprompts "github.com/FlagrantGarden/flfa/pkg/flfa/tui/traits/prompts"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func SelectProfile(settings *terminal.Settings, profiles []data.Profile) *selection.Selection {
	var profileNames []string
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name())
	}

	return selector.NewStringSelector(
		"What profile should this Group have?",
		profileNames,
		selector.WithPageSize(5),
		selector.WithSelectedChoiceStyle(selector.ColorizedBasicSelectedChoiceStyle(settings.ExtraColor("highlight"))),
	)
}

func SelectProfileModel(settings *terminal.Settings, profiles []data.Profile) *selection.Model {
	return selection.NewModel(SelectProfile(settings, profiles))
}

func GetGroupName(settings *terminal.Settings) *textinput.TextInput {
	return texter.New(
		"What should this Group be named?",
		texter.WithPlaceholder("Group Name cannot be empty"),
		texter.WithInputWidth(30),
	)
}

func GetGroupNameModel(settings *terminal.Settings) *textinput.Model {
	return textinput.NewModel(GetGroupName(settings))
}

func SelectGroupEditingOption(settings *terminal.Settings, hasAddableTraits bool, hasRemovableTraits bool) *selection.Selection {
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
		selector.WithSelectedChoiceStyle(selector.ColorizedBasicSelectedChoiceStyle(settings.ExtraColor("highlight"))),
	)
}

func SelectGroupEditingOptionModel(settings *terminal.Settings, hasAddableTraits bool, hasRemovableTraits bool) *selection.Model {
	return selection.NewModel(SelectGroupEditingOption(settings, hasAddableTraits, hasRemovableTraits))
}

func SelectAddSpecialTrait(settings *terminal.Settings, applicableTraits []data.Trait) *selection.Selection {
	add := settings.ApplyAndRender(
		"add",
		terminal.OverrideWithExtraStyle("strong"),
		terminal.ColorizeForeground("highlight"),
	)
	return tprompts.SelectTrait(
		settings,
		fmt.Sprintf("Which trait do you want to %s?", add),
		applicableTraits,
	)
}

func SelectAddSpecialTraitModel(settings *terminal.Settings, applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectAddSpecialTrait(settings, applicableTraits))
}

func SelectRemoveSpecialTrait(settings *terminal.Settings, applicableTraits []data.Trait) *selection.Selection {
	remove := settings.ApplyAndRender(
		"remove",
		terminal.OverrideWithExtraStyle("strong"),
		terminal.ColorizeForeground("warning"),
	)

	promptSettings, _ := settings.Copy()
	promptSettings.SetFlagOn("removing_trait")

	return tprompts.SelectTrait(
		promptSettings,
		fmt.Sprintf("Which trait do you want to %s?", remove),
		applicableTraits,
	)
}

func SelectRemoveSpecialTraitModel(settings *terminal.Settings, applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectRemoveSpecialTrait(settings, applicableTraits))
}

func ConfirmChangeBaseProfile(settings *terminal.Settings, new_profile string) *confirmation.Confirmation {
	var messageBuilder strings.Builder
	new_profile = settings.RenderWithDynamicStyle("confirmation_emphasis", new_profile)
	reset := settings.RenderWithDynamicStyle("warning_emphasis", "reset")
	lose := settings.RenderWithDynamicStyle("warning_emphasis", "lose")
	messageBuilder.WriteString(fmt.Sprintf("Are you sure you want to update the base profile to '%s'?\n", new_profile))
	messageBuilder.WriteString(fmt.Sprintf("Doing so will %s this group to that base profile. ", reset))
	messageBuilder.WriteString(fmt.Sprintf("You'll %s any changes other than the group's name and captaincy.", lose))
	message := lipgloss.NewStyle().Width(120).Render(messageBuilder.String())
	message = strings.Join([]string{message, ""}, "\n")

	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmChangeBaseProfileModel(settings *terminal.Settings, new_profile string) *confirmation.Model {
	return confirmation.NewModel(ConfirmChangeBaseProfile(settings, new_profile))
}
