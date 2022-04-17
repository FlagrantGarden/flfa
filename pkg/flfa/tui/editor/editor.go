package editor

import (
	"time"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	persona "github.com/FlagrantGarden/flfa/pkg/flfa/state/player"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/player"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tui.SharedModel
	Player   *player.Model
	Company  *company.Model
	Indexes  Indexes
	Substate Substate
}

const (
	StateEditingMenu compositor.State = iota + 1000
	StatePlayerMenu
	StateCompanyMenu
	StateRosterMenu
)

type Indexes struct {
	EditingCompany  int
	RemovingCompany int
}

type Substate struct {
	Editing SubstateEditing
	Player  SubstatePlayer
	Company SubstateCompany
}

func (model *Model) SetAndStartState(state compositor.State) (cmd tea.Cmd) {
	switch state {
	case StatePlayerMenu:
		cmd = model.SetAndStartSubstate(SelectingPlayer)
	case StateCompanyMenu:
		cmd = model.SetAndStartSubstate(EditingCompany)
	case StateRosterMenu:
	case StateEditingMenu:
		cmd = model.SetAndStartSubstate(SelectingOption)
	}
	return cmd
}

func (model *Model) SetAndStartSubstate(substate compositor.SubstateInterface[*Model]) (cmd tea.Cmd) {
	switch substate.(type) {
	case SubstateEditing:
		model.State = StateEditingMenu
		model.Substate.Editing = substate.(SubstateEditing)
		cmd = model.Substate.Editing.Start(model)
	case SubstatePlayer:
		model.State = StatePlayerMenu
		model.Substate.Player = substate.(SubstatePlayer)
		cmd = model.Substate.Player.Start(model)
	case SubstateCompany:
		model.State = StateCompanyMenu
		model.Substate.Company = substate.(SubstateCompany)
		cmd = model.Substate.Company.Start(model)
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

func WithPlayer(activePlayer *persona.Player) compositor.Option[*Model] {
	return func(model *Model) {
		model.Player = player.NewModel(model.Api, player.WithPlayer(activePlayer))
	}
}

func WithCompany(activeCompany *data.Company) compositor.Option[*Model] {
	return func(model *Model) {
		model.Company = company.NewModel(model.Api, company.WithCompany(activeCompany))
	}
}

func (model *Model) Init() tea.Cmd {
	if model.Player == nil {
		return model.SetAndStartState(StatePlayerMenu)
	}

	if model.Company == nil {
		return model.SetAndStartState(StateCompanyMenu)
	}

	return model.SetAndStartState(StateEditingMenu)
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
	case StateEditingMenu:
		view = model.Substate.Editing.View(model)
	case StatePlayerMenu:
		view = model.Substate.Player.View(model)
	case StateCompanyMenu:
		view = model.Substate.Company.View(model)
	case StateRosterMenu:
	case compositor.StateBroken:
		view = model.ViewFatalError()
	}

	return view
}
