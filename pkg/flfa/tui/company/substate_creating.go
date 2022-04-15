package company

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company/prompts"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubstateCreating int

const (
	Naming SubstateCreating = iota
	Describing
	AddingFirstGroup
)

func (state SubstateCreating) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		model.TextInput = prompts.GetNameModel()
		cmd = model.TextInput.Init()
	case Describing:
		model.TextInput = prompts.GetDescriptionModel()
		cmd = model.TextInput.Init()
	case AddingFirstGroup:
		model.Group = group.NewModel(model.Api, group.AsSubModel())
		cmd = model.Group.Init()
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEnded(model *Model) (cmd tea.Cmd) {
	switch state {
	case AddingFirstGroup:
		switch model.Group.State {
		case compositor.StateCancelled:
			cmd = model.Cancelled
		case compositor.StateDone:
			cmd = model.AddGroup()
		}
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		cmd = model.UpdateName(Describing)
	case Describing:
		cmd = model.UpdateDescription(AddingFirstGroup)
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
	case Describing:
	case AddingFirstGroup:
	}

	return cmd
}

func (state SubstateCreating) UpdateOnFallThrough(model *Model, msg tea.Msg) (cmd tea.Cmd) {
	switch model.Substate.Creating {
	case Naming, Describing:
		_, cmd = model.TextInput.Update(msg)
	case AddingFirstGroup:
		_, cmd = model.Group.Update(msg)
	}

	return cmd
}

func (state SubstateCreating) View(model *Model) string {
	var header string
	var subview string
	switch state {
	case Naming:
		header = "Creating a new company\n\n"
		subview = model.TextInput.View()
	case Describing:
		header = fmt.Sprintf("Creating the new %s company\n\n", model.FormattedCompanyName())
		subview = model.TextInput.View()
	case AddingFirstGroup:
		header = fmt.Sprintf(
			"Adding the first group to '%s'\n%s\n\n",
			model.FormattedCompanyName(),
			model.FormattedCompanyDescription(),
		)
		subview = model.Group.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subview,
	)
}
