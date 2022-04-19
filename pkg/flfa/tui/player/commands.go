package player

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/player"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) LoadPlayer() tea.Cmd {
	var name string
	// persona not yet set
	if model.Player == nil {
		if model.Options.LoadActivePlayer {
			name = model.Api.Tympan.Configuration.ActiveUserPersona
		}
		if name == "" {
			// no preferred user account; go choose one
			return model.SetAndStartSubstate(SelectingPersona)
		}
	} else {
		name = model.Name
	}

	foundPlayer, err := model.Api.GetPlayer(name, "")
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Player = &foundPlayer
	return model.SetAndStartState(StateEditingPersona)
}

func (model *Model) InitializePlayer(name string, nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	kind := player.Kind()
	model.Player = &player.Player{
		Persona: &persona.Persona[player.Data, player.Settings]{Kind: *kind},
	}
	err := model.Player.Initialize(name,
		model.Api.Tympan.Configuration.FolderPaths.Cache,
		model.Api.Tympan.AFS,
	)
	if err != nil {
		return model.RecordFatalError(err)
	}

	return model.SetAndStartSubstate(nextSubstate)
}

func (model *Model) SavePlayer(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	return func() tea.Msg {
		err := model.Player.Save(model.Api.Tympan.AFS)
		if err != nil {
			return model.RecordFatalError(err)
		}
		return model.SetAndStartSubstate(nextSubstate)
	}
}

func (model *Model) UpdateFallThrough(msg tea.Msg) (cmd tea.Cmd) {
	switch model.State {
	case StateChoosingPersona:
		cmd = model.Substate.Choosing.UpdateOnFallThrough(model, msg)
	case StateCreatingPersona:
		cmd = model.Substate.Creating.UpdateOnFallThrough(model, msg)
	case StateEditingPersona:
		cmd = model.Substate.Editing.UpdateOnFallThrough(model, msg)
	}

	return cmd
}

func (model *Model) UpdateOnKeyPress(msg tea.KeyMsg) (cmd tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		cmd = tea.Quit
	case "esc":
		switch model.State {
		case StateChoosingPersona:
			cmd = model.Substate.Choosing.UpdateOnEsc(model)
		case StateCreatingPersona:
			cmd = model.Substate.Creating.UpdateOnEsc(model)
		case StateEditingPersona:
			cmd = model.Substate.Editing.UpdateOnEsc(model)
		}
	case "enter":
		switch model.State {
		case StateChoosingPersona:
			cmd = model.Substate.Choosing.UpdateOnEnter(model)
		case StateCreatingPersona:
			cmd = model.Substate.Creating.UpdateOnEnter(model)
		case StateEditingPersona:
			cmd = model.Substate.Editing.UpdateOnEnter(model)
		}
	}
	return cmd
}

func (model *Model) UpdateOnSubmodelEnded() (cmd tea.Cmd) {
	// No submodels send an end message.
	return cmd
}

func (model *Model) UpdateSelectingPersona() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if choice.Value.(string) == "Create a new persona" {
		cmd = model.SetAndStartSubstate(Naming)
	} else {
		cmd = model.InitializePlayer(choice.Value.(string), SelectingEditingOption)
	}

	return cmd
}

func (model *Model) UpdateConfirmPreferred(nextState compositor.State) (cmd tea.Cmd) {
	makePreferred, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if makePreferred {
		model.Api.Tympan.Configuration.ActiveUserPersona = model.Name
		model.State = compositor.StateSavingConfiguration
		cmd = model.SaveConfig()
		if cmd == nil {
			cmd = model.SetAndStartState(nextState)
		}
	} else {
		cmd = model.SetAndStartState(nextState)
	}

	return cmd
}

func (model *Model) UpdateName(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	name, err := model.TextInput.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	return model.InitializePlayer(name, nextSubstate)
}
