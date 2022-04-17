package editor

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstatePlayer int

const (
	SelectingPlayer SubstatePlayer = iota
)

func (state SubstatePlayer) Start(model *Model) (cmd tea.Cmd) {
	model.Api.CachePlayers("")
	model.Player = player.NewModel(model.Api, player.AsSubModel())
	cmd = model.Player.Init()

	return cmd
}

func (state SubstatePlayer) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	return cmd
}

func (state SubstatePlayer) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	return cmd
}

func (state SubstatePlayer) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	switch model.Player.State {
	case compositor.StateCancelled:
		if model.Player.Player == nil {
			cmd = model.Cancelled
		} else {
			cmd = model.SetAndStartSubstate(SelectingOption)
		}
	case compositor.StateDone:
		cmd = model.SetAndStartSubstate(SelectingOption)
	}
	return cmd
}

func (state SubstatePlayer) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	_, cmd = model.Player.Update(msg)

	return cmd
}

func (state SubstatePlayer) View(model *Model) (view string) {
	return model.Player.View()
}
