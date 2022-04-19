package editor

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type SubstatePlayer int

const (
	LoadingPlayer SubstatePlayer = iota
	SelectingPlayer
)

func (state SubstatePlayer) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case LoadingPlayer:
		model.Player = player.NewModel(model.Api, player.AsSubModel(), player.WithLoadActivePlayer())
	case SelectingPlayer:
		model.Api.CachePlayers("")
		if model.Player != nil {
			model.Cache.Player = model.Player.Player
		}
		model.Player = player.NewModel(model.Api, player.AsSubModel())
	}

	return model.Player.Init()
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
			if model.Cache.Player == nil {
				cmd = model.Cancelled
			} else {
				model.Player.Player = model.Cache.Player
				cmd = model.SetAndStartSubstate(SelectingOption)
			}
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
