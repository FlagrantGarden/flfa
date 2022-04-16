package company

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company/prompts"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) UpdateFallThrough(msg tea.Msg) (cmd tea.Cmd) {
	switch model.State {
	case StateChoosingCompany:
		cmd = model.Substate.Choosing.UpdateOnFallThrough(model, msg)
	case StateCreatingCompany:
		cmd = model.Substate.Creating.UpdateOnFallThrough(model, msg)
	case StateEditingCompany:
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
		case StateChoosingCompany:
			cmd = model.Substate.Choosing.UpdateOnEnter(model)
		case StateCreatingCompany:
			cmd = model.Substate.Creating.UpdateOnEnter(model)
		case StateEditingCompany:
			cmd = model.Substate.Editing.UpdateOnEnter(model)
		}
	}

	return cmd
}

func (model *Model) UpdateOnSubmodelEnded() (cmd tea.Cmd) {
	switch model.State {
	case StateCreatingCompany:
		cmd = model.Substate.Creating.UpdateOnEnded(model)
	case StateEditingCompany:
		cmd = model.Substate.Editing.UpdateOnEnded(model)
	}
	return cmd
}

// Company Update Commands

func (model *Model) UpdateSelectingCompany() tea.Cmd {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if choice.String == "Create a new company" {
		return model.SetAndStartSubstate(Naming)
	}

	for _, company := range model.Api.Cache.Companies {
		if company.Name == choice.String {
			copyOfCompany := company
			model.Company = &copyOfCompany
			break
		}
	}

	if model.Company == nil {
		return model.RecordFatalError(fmt.Errorf("Unable to set company '%s'!", choice.String))
	}

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) AddGroup() tea.Cmd {
	model.Groups = append(model.Groups, *model.Group.Group)

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateName(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	name, err := model.TextInput.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Name = name

	return model.SetAndStartSubstate(nextSubstate)
}

func (model *Model) UpdateDescription(nextSubstate compositor.SubstateInterface[*Model]) tea.Cmd {
	description, err := model.TextInput.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Description = description

	return model.SetAndStartSubstate(nextSubstate)
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
	case "Change Description":
		cmd = model.SetAndStartSubstate(Redescribing)
	case "Change Name":
		cmd = model.SetAndStartSubstate(Renaming)
	case "Create & add a new Group":
		cmd = model.SetAndStartSubstate(AddingGroup)
	case "Edit a Group":
		cmd = model.SetAndStartSubstate(SelectingGroupToEdit)
	case "Add a copy of a Group":
		cmd = model.SetAndStartSubstate(CopyingGroup)
	case "Remove a Group":
		cmd = model.SetAndStartSubstate(RemovingGroup)
	case "Update Captaincy":
		cmd = model.SetAndStartSubstate(SelectingCaptainOption)
	case "Promote a Group to Captain":
		cmd = model.SetAndStartSubstate(SelectingGroupToPromote)
	}

	return cmd
}

func (model *Model) UpdateGroupSelection() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	selectedGroup := choice.Value.(data.Group)

	switch model.Substate.Editing {
	case SelectingGroupToEdit:
		model.Substate.Editing = EditingGroup
		model.Indexes.EditingGroup = choice.Index
		model.Group = group.NewModel(
			model.Api,
			group.AsSubModel(),
			group.WithGroup(&selectedGroup),
			group.WithCompany(model.Company),
		)

		cmd = model.Group.Init()
	case SelectingGroupToPromote:
		model.Groups[choice.Index].PromoteToCaptain(nil, data.FilterTraitsBySource("core", model.Api.Cache.Traits)...)
		cmd = model.SetAndStartSubstate(SelectingOption)
	case CopyingGroup:
		model.Groups = append(model.Groups, selectedGroup)
		cmd = model.SetAndStartSubstate(SelectingOption)
	case RemovingGroup:
		model.Groups = utils.RemoveIndex(model.Groups, choice.Index)
		cmd = model.SetAndStartSubstate(SelectingOption)
	case SelectingCaptainReplacement:
		if model.Indexes.CurrentCaptain != choice.Index {
			model.Indexes.ReplacementCaptain = choice.Index
			cmd = model.SetAndStartSubstate(ConfirmingCaptainReplacement)
			break
		}
		cmd = model.SetAndStartSubstate(SelectingOption)
	}

	return cmd
}

func (model *Model) UpdateGroup() tea.Cmd {
	model.Groups[model.Indexes.EditingGroup] = *model.Group.Group

	return model.SetAndStartSubstate(SelectingOption)
}

// Captaincy Update Commands

func (model *Model) UpdateCaptainSelection() (cmd tea.Cmd) {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	switch choice.String {
	case "Go back":
		cmd = model.SetAndStartSubstate(SelectingOption)
	case "Reroll Captain's trait":
		cmd = model.SetAndStartSubstate(RerollingCaptainTrait)
	case "Choose Captain's trait":
		cmd = model.SetAndStartSubstate(SelectingCaptainTrait)
	case "Demote Captain":
		cmd = model.SetAndStartSubstate(ConfirmingCaptainDemotion)
	case "Choose a different Captain":
		cmd = model.SetAndStartSubstate(SelectingCaptainReplacement)
	}

	return cmd
}

func (model *Model) UpdateCaptainTrait() tea.Cmd {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.Groups[model.Indexes.CurrentCaptain].Captain = choice.Value.(data.Trait)

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateRerolledCaptainTrait() tea.Cmd {
	choice, err := model.Selection.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	trait := choice.Value.(prompts.CaptainRerollChoice).Trait

	if trait.Name == "" {
		return model.SetAndStartSubstate(RerollingCaptainTrait)
	}

	model.Groups[model.Indexes.CurrentCaptain].Captain = trait

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateCaptainDemotion() tea.Cmd {
	demote, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if demote {
		model.Groups[model.Indexes.CurrentCaptain].DemoteFromCaptain()
	}

	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) UpdateCaptainReplacement() tea.Cmd {
	replace, err := model.Confirmation.Value()
	if err != nil {
		return model.RecordFatalError(err)
	}

	if replace {
		model.Groups[model.Indexes.CurrentCaptain].DemoteFromCaptain()
		model.Groups[model.Indexes.ReplacementCaptain].PromoteToCaptain(
			nil,
			data.FilterTraitsBySource("core", model.Api.Cache.Traits)...,
		)
	}

	return model.SetAndStartSubstate(SelectingOption)
}
