package editor

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateCompany int

const (
	CreatingCompany SubstateCompany = iota
	EditingCompany
)

func (state SubstateCompany) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case CreatingCompany:
		model.Company = company.NewModel(model.Api, company.AsSubModel())
	case EditingCompany:
		model.Company = company.NewModel(
			model.Api,
			company.WithCompany(&model.Player.Data.Companies[model.Indexes.EditingCompany]),
			company.AsSubModel(),
		)
	}

	return model.Company.Init()
}

func (state SubstateCompany) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	// This substate always falls through to the Company model
	return cmd
}

func (state SubstateCompany) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	// This substate always falls through to the Company model
	return cmd
}

func (state SubstateCompany) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	switch state {
	case CreatingCompany:
		switch model.Company.State {
		case compositor.StateCancelled:
			cmd = model.SetAndStartSubstate(SelectingOption)
		case compositor.StateDone:
			cmd = model.AddCompany()
		}
	case EditingCompany:
		switch model.Company.State {
		case compositor.StateCancelled:
			cmd = model.SetAndStartSubstate(SelectingOption)
		case compositor.StateDone:
			cmd = model.UpdateCompany()
		}
	}
	return cmd
}

func (state SubstateCompany) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	_, cmd = model.Company.Update(msg)

	return cmd
}

func (state SubstateCompany) View(model *Model) (view string) {
	return model.Company.View()
}
