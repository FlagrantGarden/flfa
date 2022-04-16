package company

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company/prompts"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubstateEditing int

const (
	SelectingOption SubstateEditing = iota
	Renaming
	Redescribing
	AddingGroup
	SelectingGroupToEdit
	SelectingGroupToPromote
	EditingGroup
	CopyingGroup
	RemovingGroup
	SelectingCaptainOption
	RerollingCaptainTrait
	SelectingCaptainTrait
	SelectingCaptainReplacement
	ConfirmingCaptainReplacement
	ConfirmingCaptainDemotion
)

func (state SubstateEditing) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		canRemoveAGroup := len(model.Groups) > 1
		hasCaptain := model.HasCaptain()
		model.Selection = prompts.SelectOptionModel(canRemoveAGroup, hasCaptain)
		cmd = model.Selection.Init()
	case Renaming:
		model.TextInput = prompts.GetNameModel()
		cmd = model.TextInput.Init()
	case Redescribing:
		model.TextInput = prompts.GetDescriptionModel()
		cmd = model.TextInput.Init()
	case AddingGroup:
		model.Group = group.NewModel(model.Api, group.AsSubModel(), group.WithCompany(model.Company))
		cmd = model.Group.Init()
	case SelectingGroupToEdit:
		model.Selection = prompts.SelectGroupModel(prompts.Editing, model.Groups)
		cmd = model.Selection.Init()
	case SelectingGroupToPromote:
		model.Selection = prompts.SelectGroupModel(prompts.Promoting, model.Groups)
		cmd = model.Selection.Init()
	case EditingGroup:
		// ??
	case CopyingGroup:
		model.Selection = prompts.SelectGroupModel(prompts.Copying, model.Groups)
		cmd = model.Selection.Init()
	case RemovingGroup:
		model.Selection = prompts.SelectGroupModel(prompts.Removing, model.Groups)
		cmd = model.Selection.Init()
	case SelectingCaptainOption:
		model.Selection = prompts.SelectCaptaincyOptionModel()
		cmd = model.Selection.Init()
	case RerollingCaptainTrait:
		model.Selection = prompts.SelectRerollCaptainTraitModel(
			model.CaptainsGroup(),
			data.FilterTraitsBySource("core", model.Api.Cache.Traits),
		)
		cmd = model.Selection.Init()
	case SelectingCaptainTrait:
		model.Selection = prompts.SelectCaptainTraitModel(data.FilterTraitsBySource("core", model.Api.Cache.Traits))
		cmd = model.Selection.Init()
	case SelectingCaptainReplacement:
		// promotableGroups := utils.RemoveIndex(model.Groups, model.CurrentCaptainIndex)
		model.Selection = prompts.SelectGroupModel(prompts.Promoting, model.Groups)
		cmd = model.Selection.Init()
	case ConfirmingCaptainReplacement:
		model.Confirmation = prompts.ConfirmReplaceCaptainModel(
			model.CaptainsGroup(),
			model.Groups[model.Indexes.ReplacementCaptain],
		)
		cmd = model.Confirmation.Init()
	case ConfirmingCaptainDemotion:
		model.Confirmation = prompts.ConfirmDemoteCaptainModel(model.CaptainsGroup())
		cmd = model.Confirmation.Init()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	switch state {
	case AddingGroup:
		switch model.Group.State {
		case compositor.StateCancelled:
			cmd = model.SetAndStartSubstate(SelectingOption)
		case compositor.StateDone:
			cmd = model.AddGroup()
		}
	case EditingGroup:
		switch model.Group.State {
		case compositor.StateCancelled:
			cmd = model.SetAndStartSubstate(SelectingOption)
		case compositor.StateDone:
			cmd = model.UpdateGroup()
		}
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
		return model.UpdateEditingChoice()
	case Renaming:
		return model.UpdateName(SelectingOption)
	case Redescribing:
		return model.UpdateDescription(SelectingOption)
	case SelectingGroupToEdit, SelectingGroupToPromote, RemovingGroup, SelectingCaptainReplacement, CopyingGroup:
		return model.UpdateGroupSelection()
	case SelectingCaptainOption:
		return model.UpdateCaptainSelection()
	case SelectingCaptainTrait:
		return model.UpdateCaptainTrait()
	case RerollingCaptainTrait:
		return model.UpdateRerolledCaptainTrait()
	case ConfirmingCaptainDemotion:
		return model.UpdateCaptainDemotion()
	case ConfirmingCaptainReplacement:
		return model.UpdateCaptainReplacement()
	}

	return cmd
}

func (state SubstateEditing) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	switch state {
	case SelectingOption:
	case Renaming:
	case Redescribing:
	case AddingGroup:
	case SelectingGroupToEdit:
	case SelectingGroupToPromote:
	case EditingGroup:
	case RemovingGroup:
	case SelectingCaptainOption:
	case RerollingCaptainTrait:
	case SelectingCaptainTrait:
	case SelectingCaptainReplacement:
	case ConfirmingCaptainReplacement:
	case ConfirmingCaptainDemotion:
	}

	return cmd
}

func (state SubstateEditing) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch model.Substate.Editing {
	case ConfirmingCaptainDemotion, ConfirmingCaptainReplacement:
		_, cmd = model.Confirmation.Update(msg)
	case SelectingOption, SelectingCaptainOption, RemovingGroup, CopyingGroup,
		SelectingGroupToEdit, SelectingGroupToPromote, SelectingCaptainTrait,
		RerollingCaptainTrait, SelectingCaptainReplacement:
		_, cmd = model.Selection.Update(msg)
	case Renaming, Redescribing:
		_, cmd = model.TextInput.Update(msg)
	case AddingGroup, EditingGroup:
		_, cmd = model.Group.Update(msg)
	}

	return cmd
}

func (state SubstateEditing) View(model *Model) string {
	var companyOverview string
	var captainSummary string
	var subview string

	if model.ShouldDisplayCompanyOverview() {
		companyOverview = model.CompanyOverview()
	} else if model.ShouldDisplayCaptainSummary() {
		captainSummary = model.CaptainSummary()
	}

	switch state {
	case ConfirmingCaptainDemotion, ConfirmingCaptainReplacement:
		subview = model.Confirmation.View()
	case Renaming, Redescribing:
		subview = model.TextInput.View()
	case AddingGroup, EditingGroup:
		subview = model.Group.View()
	default:
		subview = model.Selection.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		companyOverview,
		captainSummary,
		subview,
	)
}
