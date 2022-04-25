package player

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player/prompts"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateEditing int

const (
	IdleEditing SubstateEditing = iota
	SelectingEditingOption
	Renaming
	ConfirmingRename
	ConfirmingPreferredStatusUpdate
)

func (state SubstateEditing) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingEditingOption:
		// TODO; for now, just send a done message if submodel or quit
		if model.IsSubmodel {
			cmd = model.Done
		} else {
			cmd = tea.Quit
		}
	case Renaming:
		model.Selection = prompts.ChoosePlayerModel(model.TerminalSettings, model.Api.Cache.Players)
		cmd = model.Selection.Init()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingEditingOption:
		// TODO
	case Renaming:
		cmd = model.UpdateSelectingPersona()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	// TODO
	return cmd
}

func (state SubstateEditing) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	// no states have submodels that send an ended message
	return cmd
}

func (state SubstateEditing) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch state {
	case SelectingEditingOption:
		// TODO
	case Renaming:
		_, cmd = model.Selection.Update(msg)
	}

	return cmd
}

func (state SubstateEditing) View(model *Model) (view string) {
	switch state {
	case SelectingEditingOption:
		// TODO
	case Renaming:
		view = model.Selection.View()
	}

	return view
}
