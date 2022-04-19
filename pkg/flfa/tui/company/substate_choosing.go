package company

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company/prompts"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateChoosing int

const (
	SelectingCompany SubstateChoosing = iota
)

func (state SubstateChoosing) Start(model *Model) (cmd tea.Cmd) {
	switch model.Substate.Choosing {
	case SelectingCompany:
		model.Selection = prompts.ChooseCompanyModel(true, model.AvailableCompanies)
		cmd = model.Selection.Init()
	}

	return cmd
}

func (state SubstateChoosing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch model.Substate.Choosing {
	case SelectingCompany:
		cmd = model.UpdateSelectingCompany()
	}

	return cmd
}

func (state SubstateChoosing) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	if model.IsSubmodel {
		cmd = model.Cancelled
	} else {
		cmd = tea.Quit
	}
	return cmd
}

func (state SubstateChoosing) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	// no states have submodels that send an ended message
	return cmd
}

func (state SubstateChoosing) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch model.Substate.Choosing {
	case SelectingCompany:
		_, cmd = model.Selection.Update(msg)
	}

	return cmd
}

func (state SubstateChoosing) View(model *Model) (view string) {
	switch state {
	case SelectingCompany:
		view = model.Selection.View()
	}

	return view
}
