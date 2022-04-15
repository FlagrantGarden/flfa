package group

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// Required Modeler Methods

func (model *Model) UpdateFallThrough(msg tea.Msg) (cmd tea.Cmd) {
	switch model.State {
	case StateCreatingGroup:
		cmd = model.Substate.Creation.UpdateOnFallThrough(model, msg)
	case StateEditingGroup:
		cmd = model.Substate.Editing.UpdateOnFallThrough(model, msg)
	}
	return cmd
}

func (model *Model) UpdateOnKeyPress(msg tea.KeyMsg) (cmd tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	// TODO: Implement backing out of a menu with escape
	case "esc":
		switch model.State {
		case StateCreatingGroup:
			// confirm you want to cancel/exit model, goto ready state
		case StateEditingGroup:
			// return to selection menu; if on selection menu, confirm to exit model
		}
	case "enter":
		switch model.State {
		case StateCreatingGroup:
			cmd = model.Substate.Creation.UpdateOnEnter(model)
		case StateEditingGroup:
			cmd = model.Substate.Editing.UpdateOnEnter(model)
		}
	}

	return cmd
}

func (model *Model) UpdateOnSubmodelEnded() (cmd tea.Cmd) {
	// No submodels send an end message.
	return cmd
}

// Group-Specific Methods

func (model *Model) InitializeGroup(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	group, err := data.NewGroup(model.Name, model.ProfileName, model.ApplicableProfiles())
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Group = &group
	model.BaseProfile = group
	model.UpdateCompanyWorkingPointTotal()

	model.State = StateEditingGroup
	return model.SetAndStartSubstate(nextSubstate)
}

func (model *Model) UpdateName(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	name, err := model.TextInput.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Name = name
	return model.SetAndStartSubstate(nextSubstate)
}

func (model *Model) UpdateBaseProfile(updating bool) (cmd tea.Cmd) {
	// User had to confirm switching to a new base profile, need to either update the group or reset
	if updating && model.Substate.Editing == ConfirmingBaseProfileUpdate {
		shouldUpdate, err := model.Confirmation.Value()
		if err != nil {
			return model.RecordFatalError(err)
		}
		if shouldUpdate {
			model.State = StateInitializingGroup
			model.ProfileName = model.Temp.ProfileName
			return model.InitializeGroup(SelectingOption)
		} else {
			return model.SetAndStartSubstate(SelectingOption)
		}
	}

	// User selected a new base profile
	base_profile, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	// Group already had a base profile, need to confirm
	if updating {
		model.Temp.ProfileName = base_profile.String
		return model.SetAndStartSubstate(ConfirmingBaseProfileUpdate)
	}

	model.ProfileName = base_profile.String
	model.State = StateInitializingGroup
	return model.InitializeGroup(SelectingOption)
}

func (model *Model) UpdateEditingChoice() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	switch choice.String {
	case "Save & Continue":
		if model.IsSubmodel {
			model.State = compositor.StateReady
			return model.Done
		}

		cmd = tea.Quit
	case "Add Special Trait":
		cmd = model.SetAndStartSubstate(AddingSpecialTrait)
	case "Remove Special Trait":
		cmd = model.SetAndStartSubstate(RemovingSpecialTrait)
	case "Change Base Profile":
		cmd = model.SetAndStartSubstate(ChangingBaseProfile)
	case "Change Name":
		cmd = model.SetAndStartSubstate(Renaming)
	}

	return
}

func (model *Model) UpdateTrait() (cmd tea.Cmd) {
	switch model.Substate.Editing {
	case MakingTraitChoice:
		cmd = model.UpdateTraitChoice()
	case AddingSpecialTrait:
		cmd = model.UpdateTraitAdd()
	case RemovingSpecialTrait:
		cmd = model.UpdateTraitRemove()
	}

	return
}

func (model *Model) UpdateTraitChoice() (cmd tea.Cmd) {
	trait := model.CurrentTraitWithChoice()
	choice := model.CurrentTraitChoice()

	switch choice.Prompt.Type {
	case "text":
		chooser := model.TraitChooser.TextInput
		choice, err := chooser.Value()
		if err != nil {
			return model.RecordFatalError(err)
		}

		trait.Choices[model.Indexes.CurrentChoice].Value = choice
		model.UpdateCurrentTraitWithChoiceName()
		model.Group, err = trait.AddToGroup(model.Group, model.Api.ScriptEngine)
		if err != nil {
			return model.RecordFatalError(err)
		}

		if (model.Indexes.CurrentChoice + 1) < len(trait.Choices) {
			model.Indexes.CurrentChoice += 1
			return model.SetAndStartSubstate(MakingTraitChoice)
		}

		cmd = model.SetAndStartSubstate(SelectingOption)
	case "selection": // TODO: Implement
	case "confirmation": // TODO: Implement
	}

	return
}

func (model *Model) UpdateTraitAdd() tea.Cmd {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}
	chosenTrait := choice.Value.(data.Trait)
	if len(chosenTrait.Choices) > 0 {
		model.TraitsWithChoices = append(model.TraitsWithChoices, &chosenTrait)
		model.Indexes.CurrentTraitWithChoice = len(model.TraitsWithChoices) - 1
		model.Indexes.CurrentChoice = 0
		return model.SetAndStartSubstate(MakingTraitChoice)
	}

	updatedGroup, err := choice.Value.(data.Trait).AddToGroup(model.Group, model.Api.ScriptEngine)
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Group = updatedGroup
	model.UpdateCompanyWorkingPointTotal()

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateTraitRemove() tea.Cmd {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}
	traitToRemove := choice.Value.(data.Trait)
	updatedGroup, err := traitToRemove.RemoveFromGroup(model.Group, model.Api.ScriptEngine)
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Group = updatedGroup
	model.UpdateCompanyWorkingPointTotal()

	for index, trait := range model.TraitsWithChoices {
		if trait.Name == traitToRemove.Name {
			model.TraitsWithChoices = utils.RemoveIndex(model.TraitsWithChoices, index)
		}
	}

	return model.SetAndStartSubstate(SelectingOption)
}
