package editor

import (
	company_prompts "github.com/FlagrantGarden/flfa/pkg/flfa/tui/company/prompts"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/editor/prompts"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateEditing int

const (
	SelectingOption SubstateEditing = iota
	SelectingCompanyToEdit
	SelectingCompanyToRemove
	ConfirmRemoval
	ConfirmSave
	ConfirmQuitWithoutSaving
)

func (state SubstateEditing) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		hasCompanies := len(model.Player.Data.Companies) > 0
		model.Selection = prompts.SelectMenuOptionModel(hasCompanies)
		cmd = model.Selection.Init()
	case SelectingCompanyToEdit:
		model.Selection = company_prompts.ChooseCompanyModel(false, model.Player.Data.Companies)
		cmd = model.Selection.Init()
	case SelectingCompanyToRemove:
		model.Selection = company_prompts.ChooseCompanyModel(false, model.Player.Data.Companies)
		cmd = model.Selection.Init()
	case ConfirmRemoval:
		companyToRemove := model.Player.Data.Companies[model.Indexes.RemovingCompany].Name
		model.Confirmation = prompts.ConfirmRemoveCompanyModel(companyToRemove)
		cmd = model.Confirmation.Init()
	case ConfirmQuitWithoutSaving:
		model.Confirmation = prompts.ConfirmQuitWithoutSavingModel()
		cmd = model.Confirmation.Init()
	case ConfirmSave:
		model.Confirmation = prompts.ConfirmSavePlayerModel()
		cmd = model.Confirmation.Init()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		cmd = model.UpdateSelectMenuOption()
	case SelectingCompanyToEdit:
		cmd = model.UpdateSelectPlayerCompany()
	case SelectingCompanyToRemove:
		cmd = model.UpdateSelectPlayerCompany()
	case ConfirmRemoval:
		cmd = model.UpdateConfirmRemoval()
	case ConfirmQuitWithoutSaving:
		cmd = model.UpdateConfirmQuitWithoutSaving()
	case ConfirmSave:
		cmd = model.UpdateConfirmSave(StateEditingMenu)
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		// Prompt to save if changes made
	case SelectingCompanyToEdit:
		cmd = model.SetAndStartSubstate(SelectingOption)
	case ConfirmSave:
		cmd = model.SetAndStartSubstate(SelectingOption)
	case ConfirmQuitWithoutSaving:
		cmd = model.SetAndStartSubstate(SelectingOption)
	}
	return cmd
}

func (state SubstateEditing) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	// no states have submodels that send an ended message
	return cmd
}

func (state SubstateEditing) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch state {
	case SelectingOption, SelectingCompanyToEdit, SelectingCompanyToRemove:
		_, cmd = model.Selection.Update(msg)
	case ConfirmSave, ConfirmRemoval, ConfirmQuitWithoutSaving:
		_, cmd = model.Confirmation.Update(msg)
	}

	return cmd
}

func (state SubstateEditing) View(model *Model) (view string) {
	switch state {
	case SelectingOption, SelectingCompanyToEdit, SelectingCompanyToRemove:
		view = model.Selection.View()
	case ConfirmSave, ConfirmRemoval, ConfirmQuitWithoutSaving:
		view = model.Confirmation.View()
	}

	return view
}
