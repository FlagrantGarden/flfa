package persona

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/user"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) LoadPersona() tea.Cmd {
	return func() tea.Msg {
		var name string
		// persona not yet set
		if model.Persona == nil {
			active := model.Api.Tympan.Configuration.ActiveUserPersona
			if active != "" {
				name = active
			} else {
				// no preferred user account; go choose one
				return model.SetAndStartSubstate(SelectingPersona)
			}
		} else {
			name = model.Name
		}

		userPersona, err := model.Api.GetUserPersona(name, "")
		if err != nil {
			return model.RecordFatalError(err)
		}

		model.Persona = userPersona
		return model.SetAndStartState(StateEditingPersona)
	}
}

func (model *Model) InitializePersona(name string, nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	return func() tea.Msg {
		kind := user.Kind()
		model.Persona = &persona.Persona[user.Data, user.Settings]{Kind: *kind}
		err := model.Persona.Initialize(name,
			model.Api.Tympan.Configuration.FolderPaths.Cache,
			model.Api.Tympan.AFS,
		)
		if err != nil {
			return model.RecordFatalError(err)
		}

		return model.SetAndStartSubstate(nextSubstate)
	}
}

func (model *Model) SavePersona(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	return func() tea.Msg {
		err := model.Persona.Save(model.Api.Tympan.AFS)
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
		// TODO
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
		cmd = model.InitializePersona(choice.Value.(string), SelectingEditingOption)
	}

	return cmd
}

func (model *Model) UpdateConfirmPreferred() (cmd tea.Cmd) {
	makePreferred, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if makePreferred {
		model.Api.Tympan.Configuration.ActiveUserPersona = model.Name
		model.State = compositor.StateSavingConfiguration
		cmd = model.SaveConfig()
	}

	return cmd
}

func (model *Model) UpdateName(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	name, err := model.TextInput.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	return model.InitializePersona(name, nextSubstate)
}
