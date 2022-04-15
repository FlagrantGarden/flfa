package company

import (
	"fmt"
	"time"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/group"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tui.SharedModel
	*data.Company
	Limits   Limits
	Indexes  Indexes
	Substate Substate
	Group    *group.Model
}

type Limits struct {
	CompanyMaximumPoints int
	GroupMaximumPoints   int
}

type Indexes struct {
	EditingGroup       int
	CurrentCaptain     int
	ReplacementCaptain int
}

const (
	StateChoosingCompany compositor.State = iota + 200
	StateCreatingCompany
	StateEditingCompany
	StateLoadingCompany
	StateLoadedCompany
	StateSavingCompany
	StateSavedCompany
)

type Substate struct {
	Choosing SubstateChoosing
	Creating SubstateCreating
	Editing  SubstateEditing
}

func (model *Model) SetAndStartState(state compositor.State) (cmd tea.Cmd) {
	switch state {
	case StateChoosingCompany:
		cmd = model.SetAndStartSubstate(SelectingCompany)
	case StateCreatingCompany:
		cmd = model.SetAndStartSubstate(Naming)
	case StateEditingCompany:
		cmd = model.SetAndStartSubstate(SelectingOption)
	case StateLoadedCompany:
	case StateLoadingCompany:
	case StateSavedCompany:
	case StateSavingCompany:
	}

	return cmd
}

func (model *Model) SetAndStartSubstate(substate compositor.SubstateInterface[*Model]) (cmd tea.Cmd) {
	switch substate.(type) {
	case SubstateChoosing:
		model.State = StateChoosingCompany
		model.Substate.Choosing = substate.(SubstateChoosing)
		cmd = model.Substate.Choosing.Start(model)
	case SubstateCreating:
		model.State = StateCreatingCompany
		model.Substate.Creating = substate.(SubstateCreating)
		cmd = model.Substate.Creating.Start(model)
	case SubstateEditing:
		model.State = StateEditingCompany
		model.Substate.Editing = substate.(SubstateEditing)
		cmd = model.Substate.Editing.Start(model)
	}

	return cmd
}

func NewModel(api *flfa.Api, options ...compositor.Option[*Model]) *Model {
	model := &Model{
		SharedModel: tui.SharedModel{
			Api: api,
		},
	}

	for _, option := range options {
		option(model)
	}

	return model
}

func WithCompany(company *data.Company) compositor.Option[*Model] {
	return func(model *Model) {
		model.Company = company
	}
}

func WithCompanyMaxPoints(max int) compositor.Option[*Model] {
	return func(model *Model) {
		model.Limits.CompanyMaximumPoints = max
	}
}

func WithGroupMaxPoints(max int) compositor.Option[*Model] {
	return func(model *Model) {
		model.Limits.GroupMaximumPoints = max
	}
}

func (model *Model) Init() tea.Cmd {
	if model.Company == nil {
		model.Company = &data.Company{}
		return model.SetAndStartSubstate(SelectingCompany)
	}
	return model.SetAndStartSubstate(SelectingOption)
}

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// For some reason, a race condition on first update occurs
	// Sleeping for a few milliseconds is enough to prevent it.
	time.Sleep(time.Duration(5) * time.Millisecond)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Only some key presses are handled by this model;
		// all others fall through to the submodels
		cmd := model.UpdateOnKeyPress(msg)
		if cmd != nil {
			return model, cmd
		}
	case compositor.EndMsg:
		return model, model.UpdateOnSubmodelEnded()
	}

	// Passthru to submodel
	return model, model.UpdateFallThrough(msg)
}

func (model *Model) View() (view string) {
	switch model.State {
	case StateChoosingCompany:
		view = model.Substate.Choosing.View(model)
	case StateCreatingCompany:
		view = model.Substate.Creating.View(model)
	case StateEditingCompany:
		view = model.Substate.Editing.View(model)
	case compositor.StateBroken:
		view = fmt.Sprintf("Error! %s", model.FatalError)
	case compositor.StateReady:
		view = "Ready for anything..."
	}
	return view
}

func (model *Model) HasCaptain() bool {
	for index, group := range model.Groups {
		if group.Captain.Name != "" {
			model.Indexes.CurrentCaptain = index
			return true
		}
	}
	return false
}

func (model *Model) CaptainsGroup() (captain data.Group) {
	for _, group := range model.Groups {
		if group.Captain.Name != "" {
			captain = group
			break
		}
	}

	return captain
}
