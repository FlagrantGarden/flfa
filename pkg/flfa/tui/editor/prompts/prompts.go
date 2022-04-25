package prompts

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
)

func SelectMenuOption(settings *terminal.Settings, hasCompanies bool) *selection.Selection {
	options := []string{"Create a Company"}

	if hasCompanies {
		options = append(options, []string{"Edit a Company", "Remove a Company"}...)
	}

	options = append(options, []string{"Change Player", "Save", "Quit"}...)

	return selector.NewStringSelector(
		"What would you like to do?",
		options,
		selector.WithPageSize(5),
	)
}

func SelectMenuOptionModel(settings *terminal.Settings, hasCompanies bool) *selection.Model {
	return selection.NewModel(SelectMenuOption(settings, hasCompanies))
}

func ConfirmSavePlayer(settings *terminal.Settings) *confirmation.Confirmation {
	message := fmt.Sprintf(
		"Are you sure you want to save? this will %s your previous data and settings",
		settings.ApplyAndRender("overwrite", terminal.ColorizeForeground("warning")),
	)
	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.Yes))
}

func ConfirmSavePlayerModel(settings *terminal.Settings) *confirmation.Model {
	return confirmation.NewModel(ConfirmSavePlayer(settings))
}

func ConfirmQuitWithoutSaving(settings *terminal.Settings) *confirmation.Confirmation {
	message := fmt.Sprintf(
		"Do you want to save before you quit? You have %s if you don't.",
		settings.ApplyAndRender("unsaved changes that will be lost", terminal.ColorizeForeground("warning")),
	)
	return confirmer.New(
		message,
		confirmer.WithDefaultValue(confirmation.Yes),
	)
}

func ConfirmQuitWithoutSavingModel(settings *terminal.Settings) *confirmation.Model {
	return confirmation.NewModel(ConfirmQuitWithoutSaving(settings))
}

func ConfirmRemoveCompany(settings *terminal.Settings, name string) *confirmation.Confirmation {
	// company := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("166")).Render(name)
	company := settings.RenderWithDynamicStyle("warning_emphasis", name)
	return confirmer.New(
		fmt.Sprintf("Are you sure you want to remove the %s Company? This can't be undone.", company),
		confirmer.WithDefaultValue(confirmation.No),
		confirmer.WithInvertedColorTemplate(),
	)
}

func ConfirmRemoveCompanyModel(settings *terminal.Settings, name string) *confirmation.Model {
	return confirmation.NewModel(ConfirmRemoveCompany(settings, name))
}
