package prompts

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
)

func SelectMenuOption(hasCompanies bool) *selection.Selection {
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

func SelectMenuOptionModel(hasCompanies bool) *selection.Model {
	return selection.NewModel(SelectMenuOption(hasCompanies))
}

func ConfirmSavePlayer() *confirmation.Confirmation {
	return confirmer.New(
		"Are you sure you want to save? This will overwrite your previous data and settings.",
		confirmer.WithDefaultValue(confirmation.Yes),
	)
}

func ConfirmSavePlayerModel() *confirmation.Model {
	return confirmation.NewModel(ConfirmSavePlayer())
}

func ConfirmQuitWithoutSaving() *confirmation.Confirmation {
	return confirmer.New(
		"Do you want to save before you quit? You have unsaved changes that will be lost if you don't.",
		confirmer.WithDefaultValue(confirmation.Yes),
	)
}

func ConfirmQuitWithoutSavingModel() *confirmation.Model {
	return confirmation.NewModel(ConfirmQuitWithoutSaving())
}

func ConfirmRemoveCompany(name string) *confirmation.Confirmation {
	company := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("166")).Render(name)
	return confirmer.New(
		fmt.Sprintf("Are you sure you want to remove the %s Company? This can't be undone.", company),
		confirmer.WithDefaultValue(confirmation.No),
		confirmer.WithInvertedColorTemplate(),
	)
}

func ConfirmRemoveCompanyModel(name string) *confirmation.Model {
	return confirmation.NewModel(ConfirmRemoveCompany(name))
}
