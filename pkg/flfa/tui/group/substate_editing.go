package group

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group/prompts"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/dynamic"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubstateEditing int

const (
	IdleEditing SubstateEditing = iota
	SelectingOption
	Renaming
	ChangingBaseProfile
	ConfirmingBaseProfileUpdate
	AddingSpecialTrait
	RemovingSpecialTrait
	MakingTraitChoice
)

func (state SubstateEditing) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		canAddTraits := len(model.ApplicableTraits()) > 0
		canRemoveTraits := len(model.RemovableTraits()) > 0
		model.Selection = prompts.SelectGroupEditingOptionModel(model.TerminalSettings, canAddTraits, canRemoveTraits)
		cmd = model.Selection.Init()
	case Renaming:
		model.TextInput = prompts.GetGroupNameModel(model.TerminalSettings)
		cmd = model.TextInput.Init()
	case ChangingBaseProfile:
		model.Selection = prompts.SelectProfileModel(model.TerminalSettings, model.ApplicableProfiles())
		cmd = model.Selection.Init()
	case ConfirmingBaseProfileUpdate:
		model.Confirmation = prompts.ConfirmChangeBaseProfileModel(model.TerminalSettings, model.Temp.ProfileName)
		cmd = model.Confirmation.Init()
	case AddingSpecialTrait:
		model.Selection = prompts.SelectAddSpecialTraitModel(model.TerminalSettings, model.ApplicableTraits())
		cmd = model.Selection.Init()
	case RemovingSpecialTrait:
		model.Selection = prompts.SelectRemoveSpecialTraitModel(model.TerminalSettings, model.RemovableTraits())
		cmd = model.Selection.Init()
	case MakingTraitChoice:
		model.TraitChooser = dynamic.New(model.CurrentTraitChoice().Prompt)
		cmd = model.TraitChooser.Init()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		cmd = model.UpdateEditingChoice()
	case Renaming:
		cmd = model.UpdateName(SelectingOption)
	case ChangingBaseProfile, ConfirmingBaseProfileUpdate:
		cmd = model.UpdateBaseProfile(true)
	case MakingTraitChoice, AddingSpecialTrait, RemovingSpecialTrait:
		cmd = model.UpdateTrait()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		if model.IsSubmodel {
			cmd = model.Cancelled
		} else {
			cmd = tea.Quit
		}
	case Renaming, ChangingBaseProfile, AddingSpecialTrait, RemovingSpecialTrait, MakingTraitChoice:
		cmd = model.SetAndStartSubstate(SelectingOption)
	case ConfirmingBaseProfileUpdate:
		cmd = model.SetAndStartSubstate(ChangingBaseProfile)
	}
	return cmd
}

func (state SubstateEditing) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	// no states have submodels that send an ended message
	return cmd
}

func (state SubstateEditing) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch state {
	case SelectingOption, ChangingBaseProfile, AddingSpecialTrait, RemovingSpecialTrait:
		_, cmd = model.Selection.Update(msg)
	case Renaming:
		_, cmd = model.TextInput.Update(msg)
	case ConfirmingBaseProfileUpdate:
		_, cmd = model.Confirmation.Update(msg)
	case MakingTraitChoice:
		_, cmd = model.TraitChooser.Update(msg)
	}

	return cmd
}

func (state SubstateEditing) View(model *Model) (view string) {
	header := model.GroupEditingOverview()
	var subview string
	switch state {
	case SelectingOption, AddingSpecialTrait, RemovingSpecialTrait:
		subview = model.Selection.View()
	case Renaming:
		subview = lipgloss.JoinVertical(
			lipgloss.Left,
			"Renaming the Group:",
			model.TextInput.View(),
		)
	case ChangingBaseProfile:
		subview = lipgloss.JoinVertical(
			lipgloss.Left,
			"Modifying the Base Profile:",
			model.Selection.View(),
		)
	case ConfirmingBaseProfileUpdate:
		subview = model.Confirmation.View()
	case MakingTraitChoice:
		subview = model.TraitChooser.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subview,
	)
}
