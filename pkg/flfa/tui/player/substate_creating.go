package player

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player/prompts"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstateCreating int

const (
	Naming SubstateCreating = iota
	DecidingIfPreferred
)

func (state SubstateCreating) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		model.TextInput = prompts.GetNameModel()
		cmd = model.TextInput.Init()
	case DecidingIfPreferred:
		model.Confirmation = prompts.SetAsPreferredModel(model.TerminalSettings, model.Name)
		cmd = model.Confirmation.Init()
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		cmd = model.UpdateName(DecidingIfPreferred)
	case DecidingIfPreferred:
		cmd = model.UpdateConfirmPreferred(StateEditingPersona)
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		if model.IsSubmodel {
			cmd = model.Cancelled
		} else {
			cmd = tea.Quit
		}
	case DecidingIfPreferred:
		cmd = model.UpdateConfirmPreferred(StateEditingPersona)
	}
	return cmd
}

func (state SubstateCreating) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	// no states have submodels that send an ended message
	return cmd
}

func (state SubstateCreating) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch state {
	case Naming:
		_, cmd = model.TextInput.Update(msg)
	case DecidingIfPreferred:
		_, cmd = model.Confirmation.Update(msg)
	}

	return cmd
}

func (state SubstateCreating) View(model *Model) (view string) {
	switch state {
	case Naming:
		view = model.TextInput.View()
	case DecidingIfPreferred:
		view = model.Confirmation.View()
	}

	return view
}
