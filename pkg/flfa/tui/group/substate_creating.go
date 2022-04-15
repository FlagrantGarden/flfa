package group

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group/prompts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubstateCreating int

const (
	IdleCreating SubstateCreating = iota
	Naming
	SelectingProfile
)

func (state SubstateCreating) Start(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		model.TextInput = prompts.GetGroupNameModel()
		cmd = model.TextInput.Init()
	case SelectingProfile:
		model.Selection = prompts.SelectProfileModel(model.ApplicableProfiles())
		cmd = model.Selection.Init()
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEnter(model *Model) (cmd tea.Cmd) {
	switch state {
	case Naming:
		cmd = model.UpdateName(SelectingProfile)
	case SelectingProfile:
		cmd = model.UpdateBaseProfile(false)
	}

	return cmd
}

func (state SubstateCreating) UpdateOnEsc(model *Model) (cmd tea.Cmd) {
	// TODO
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
	case SelectingProfile:
		_, cmd = model.Selection.Update(msg)
	}

	return cmd
}

func (state SubstateCreating) View(model *Model) (view string) {
	var header string
	var subview string
	switch state {
	case Naming:
		header = "Creating a new Group"
		subview = model.TextInput.View()
	case SelectingProfile:
		header = fmt.Sprintf("Creating the %s Group", model.FormattedGroupName())
		subview = model.Selection.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subview,
	)
}
