package persona

import (
	"fmt"
	"time"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/user"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tui.SharedModel
	*persona.Persona[user.Data, user.Settings]
	Substate Substate
	// Edit
}

const (
	StateChoosingPersona compositor.State = iota + 100
	StateCreatingPersona
	StateEditingPersona
)

type Substate struct {
	Choosing SubstateChoosing
	Creating SubstateCreating
	Editing  SubstateEditing
}

func (model *Model) SetAndStartState(state compositor.State) (cmd tea.Cmd) {
	switch state {
	case StateChoosingPersona:
		model.SetAndStartSubstate(SelectingPersona)
	case StateCreatingPersona:
		model.SetAndStartSubstate(Naming)
	case StateEditingPersona:
		model.SetAndStartSubstate(SelectingEditingOption)
	case compositor.StateReady:
		model.State = compositor.StateReady
		cmd = nil
	}

	return cmd
}

func (model *Model) SetAndStartSubstate(substate compositor.SubstateInterface[*Model]) (cmd tea.Cmd) {
	switch substate.(type) {
	case SubstateChoosing:
		model.State = StateChoosingPersona
		model.Substate.Choosing = substate.(SubstateChoosing)
		cmd = model.Substate.Choosing.Start(model)
	case SubstateCreating:
		model.State = StateCreatingPersona
		model.Substate.Creating = substate.(SubstateCreating)
		cmd = model.Substate.Creating.Start(model)
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

func WithPersona(persona *persona.Persona[user.Data, user.Settings]) compositor.Option[*Model] {
	return func(model *Model) {
		model.Persona = persona
	}
}

func (model *Model) Init() tea.Cmd {
	return model.LoadPersona()
}

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// For some reason, a race condition on first update occurs
	// Sleeping for a few milliseconds is enough to prevent it.
	time.Sleep(time.Duration(5) * time.Millisecond)
	switch msg := msg.(type) {
	// When a key is pressed...
	case tea.KeyMsg:
		cmd := model.UpdateOnKeyPress(msg)
		if cmd != nil {
			return model, cmd
		}
	}

	// Passthru to sub-model
	return model, model.UpdateFallThrough(msg)
}

func (model *Model) View() string {
	switch model.State {
	case StateChoosingPersona:
		return model.Substate.Choosing.View(model)
	case StateCreatingPersona:
		return model.Substate.Creating.View(model)
	case StateEditingPersona:
		return model.Substate.Editing.View(model)

	case compositor.StateBroken:
		return fmt.Sprintf("Fatal Error: %s\n\nPress ctrl+c to exit.", model.FatalError)
	case compositor.StateReady:
		return fmt.Sprintf("Ready for anything!")
	}
	return ""
}
