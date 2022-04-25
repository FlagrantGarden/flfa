package editor

import (
	"reflect"

	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) UpdateFallThrough(msg tea.Msg) (cmd tea.Cmd) {
	switch model.State {
	case StatePlayerMenu:
		cmd = model.Substate.Player.UpdateOnFallThrough(model, msg)
	case StateCompanyMenu:
		cmd = model.Substate.Company.UpdateOnFallThrough(model, msg)
	case StateRosterMenu:
	case StateEditingMenu:
		cmd = model.Substate.Editing.UpdateOnFallThrough(model, msg)
	}

	return cmd
}

func (model *Model) UpdateOnKeyPress(msg tea.KeyMsg) (cmd tea.Cmd) {
	switch model.State {
	case compositor.StateBroken, compositor.StateCancelled, compositor.StateDone:
		return tea.Quit
	}

	switch msg.String() {
	case "ctrl+c":
		cmd = tea.Quit
	case "esc":
		switch model.State {
		case StateCompanyMenu:
			cmd = model.Substate.Company.UpdateOnEsc(model)
		case StateEditingMenu:
			cmd = model.Substate.Editing.UpdateOnEsc(model)
		case StatePlayerMenu:
			cmd = model.Substate.Player.UpdateOnEsc(model)
		case StateRosterMenu:
		}
	case "enter":
		switch model.State {
		case StateCompanyMenu:
			cmd = model.Substate.Company.UpdateOnEnter(model)
		case StateEditingMenu:
			cmd = model.Substate.Editing.UpdateOnEnter(model)
		case StatePlayerMenu:
			cmd = model.Substate.Player.UpdateOnEnter(model)
		case StateRosterMenu:
		}
	}

	return cmd
}

func (model *Model) UpdateOnSubmodelEnded() (cmd tea.Cmd) {
	switch model.State {
	case StatePlayerMenu:
		cmd = model.Substate.Player.UpdateOnEnded(model)
	case StateCompanyMenu:
		cmd = model.Substate.Company.UpdateOnEnded(model)
	}
	return cmd
}

func (model *Model) UpdateSelectMenuOption() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()

	if err != nil {
		return model.RecordFatalError(err)
	}

	switch choice.String {
	case "Create a Company":
		cmd = model.SetAndStartSubstate(CreatingCompany)
	case "Edit a Company":
		cmd = model.SetAndStartSubstate(SelectingCompanyToEdit)
	case "Remove a Company":
		cmd = model.SetAndStartSubstate(SelectingCompanyToRemove)
	case "Change Player":
		cmd = model.SetAndStartSubstate(SelectingPlayer)
	case "Save":
		cmd = model.SetAndStartSubstate(ConfirmSave)
	case "Quit":
		// need to make sure that the version on disk is the same as in memory;
		// if it is, quit. if it is not, prompt to save and then quit.
		savedPlayerPersona, err := model.Api.GetPlayer(model.Player.Name, "")
		if err != nil {
			model.RecordFatalError(err)
		}
		dataEqual := reflect.DeepEqual(savedPlayerPersona.Data, model.Player.Data)
		settingsEqual := reflect.DeepEqual(savedPlayerPersona.Settings, model.Player.Settings)

		if dataEqual && settingsEqual {
			cmd = model.Done
		} else {
			cmd = model.SetAndStartSubstate(ConfirmQuitWithoutSaving)
		}
	}

	return cmd
}

func (model *Model) UpdateSelectPlayerCompany() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}
	var companyIndex int
	for index, availableCompany := range model.Player.Data.Companies {
		if availableCompany.Name == choice.String {
			companyIndex = index
			break
		}
	}

	switch model.Substate.Editing {
	case SelectingCompanyToEdit:
		model.Indexes.EditingCompany = companyIndex
		cmd = model.SetAndStartSubstate(EditingCompany)
	case SelectingCompanyToRemove:
		model.Indexes.RemovingCompany = companyIndex
		cmd = model.SetAndStartSubstate(ConfirmRemoval)
	}

	return cmd
}

func (model *Model) AddCompany() (cmd tea.Cmd) {
	model.Player.Data.Companies = append(model.Player.Data.Companies, *model.Company.Company)

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateCompany() (cmd tea.Cmd) {
	model.Player.Data.Companies[model.Indexes.EditingCompany] = *model.Company.Company

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateConfirmRemoval() (cmd tea.Cmd) {
	remove, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if remove {
		model.Player.Data.Companies = utils.RemoveIndex(model.Player.Data.Companies, model.Indexes.RemovingCompany)
	}

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateConfirmSave(nextState compositor.State) (cmd tea.Cmd) {
	save, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if save {
		err = model.Player.Save(model.Api.Tympan.AFS)
		if err != nil {
			return model.RecordFatalError(err)
		}
	}

	return model.SetAndStartState(nextState)
}

func (model *Model) UpdateConfirmQuitWithoutSaving() (cmd tea.Cmd) {
	save, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if save {
		err = model.Player.Save(model.Api.Tympan.AFS)
		if err != nil {
			return model.RecordFatalError(err)
		}
	}

	return model.Done
}

func (model *Model) Quit() (cmd tea.Cmd) {
	// need to make sure that the version on disk is the same as in memory;
	// if it is, quit. if it is not, prompt to save and then quit.
	savedPlayerPersona, err := model.Api.GetPlayer(model.Player.Name, "")
	if err != nil {
		model.RecordFatalError(err)
	}
	dataEqual := reflect.DeepEqual(savedPlayerPersona.Data, model.Player.Data)
	settingsEqual := reflect.DeepEqual(savedPlayerPersona.Settings, model.Player.Settings)

	if dataEqual && settingsEqual {
		cmd = model.Done
	} else {
		cmd = model.SetAndStartSubstate(ConfirmQuitWithoutSaving)
	}

	return cmd
}
