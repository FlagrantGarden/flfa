package player

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player/prompts"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateChoosing int

const (
	SelectingPersona SubstateChoosing = iota
)

func (state SubstateChoosing) Start(model *Model) (cmd tea.Cmd) {
	switch model.Substate.Choosing {
	case SelectingPersona:
		model.Selection = prompts.ChoosePlayerModel(model.TerminalSettings, model.Api.Cache.Players)
		cmd = model.Selection.Init()
	}

	return cmd
}

func (state SubstateChoosing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch model.Substate.Choosing {
	case SelectingPersona:
		cmd = model.UpdateSelectingPersona()
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
	case SelectingPersona:
		_, cmd = model.Selection.Update(msg)
	}

	return cmd
}

func (state SubstateChoosing) View(model *Model) (view string) {
	switch state {
	case SelectingPersona:
		view = model.Selection.View()
	}

	return view
}
