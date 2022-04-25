package group

import (
	"time"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/dynamic"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tui.SharedModel
	*data.Group
	BaseProfile       data.Group
	Company           *data.Company
	Indexes           Indexes
	Limits            Limits
	Substate          Substate
	Temp              Temp
	TraitsWithChoices []*data.Trait
	// Could be any prompt, depends on the current choice
	TraitChooser *dynamic.Model
}

type Limits struct {
	CompanyMaximumPoints int
	GroupMaximumPoints   int
}

type Indexes struct {
	CurrentTraitWithChoice int
	CurrentChoice          int
}

type Temp struct {
	CompanyPoints     int
	ProfileName       string
	TraitsWithChoices []*data.Trait
}

const (
	StateCreatingGroup compositor.State = iota + 300
	StateEditingGroup
	StateInitializingGroup
	StateInitializedGroup
)

type Substate struct {
	Creation SubstateCreating
	Editing  SubstateEditing
}

func (model *Model) SetAndStartState(state compositor.State) (cmd tea.Cmd) {
	switch state {
	case StateCreatingGroup:
		cmd = model.SetAndStartSubstate(Naming)
	case StateEditingGroup:
		cmd = model.SetAndStartSubstate(SelectingOption)
	case StateInitializingGroup:
	case StateInitializedGroup:
	}

	return cmd
}

func (model *Model) SetAndStartSubstate(substate compositor.SubstateInterface[*Model]) (cmd tea.Cmd) {
	switch substate.(type) {
	case SubstateCreating:
		model.State = StateCreatingGroup
		model.Substate.Creation = substate.(SubstateCreating)
		cmd = model.Substate.Creation.Start(model)
	case SubstateEditing:
		model.State = StateEditingGroup
		model.Substate.Editing = substate.(SubstateEditing)
		cmd = model.Substate.Editing.Start(model)
	}

	return cmd
}

type Option func(model *Model)

func NewModel(api *flfa.Api, options ...compositor.Option[*Model]) *Model {
	model := &Model{
		SharedModel: tui.SharedModel{
			Api: api,
			Compositor: compositor.Compositor{
				TerminalSettings: tui.TerminalSettings(),
			},
		},
		Limits: Limits{
			CompanyMaximumPoints: 24,
			GroupMaximumPoints:   12,
		},
	}

	for _, option := range options {
		option(model)
	}

	model.UpdateCompanyWorkingPointTotal()

	return model
}

func WithGroup(group *data.Group) compositor.Option[*Model] {
	return func(model *Model) {
		model.Group = group
	}
}

func WithGroupMaxPoints(max int) compositor.Option[*Model] {
	return func(model *Model) {
		model.Limits.GroupMaximumPoints = max
	}
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

func AsSubModel() compositor.Option[*Model] {
	return func(model *Model) {
		model.IsSubmodel = true
	}
}

func (model *Model) Init() tea.Cmd {
	if model.Group == nil {
		model.Group = &data.Group{}
		return model.SetAndStartState(StateCreatingGroup)
	}
	return model.SetAndStartState(StateEditingGroup)
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
	}

	// Passthru to sub-model
	cmd := model.UpdateFallThrough(msg)
	if cmd != nil {
		return model, cmd
	}

	return model, nil
}

func (model *Model) View() (view string) {
	switch model.State {
	case StateCreatingGroup:
		view = model.Substate.Creation.View(model)
	case StateEditingGroup:
		view = model.Substate.Editing.View(model)
	case StateInitializingGroup:
		view = "Waiting for group to initialize..."
	case compositor.StateBroken:
		view = model.ViewFatalError()
	case compositor.StateReady:
		view = "Ready for anything..."
	}
	return view
}

func (model *Model) Return() *data.Group {
	return model.Group
}

func (model *Model) ApplicableTraits() (applicableTraits []data.Trait) {
	var errors []error
	traits := data.FilterTraitsByType("Special", model.Api.Cache.Traits)
	for _, trait := range traits {
		applicable, err := trait.Applicable(
			*model.Group,
			model.BaseProfile,
			model.Api.ScriptEngine,
			data.WithGroupMaxPoints(model.Limits.GroupMaximumPoints),
			data.WithCompanyMaxPoints(model.Limits.CompanyMaximumPoints, model.Temp.CompanyPoints),
		)
		if err != nil {
			errors = append(errors, err)
		}
		if applicable {
			applicableTraits = append(applicableTraits, trait)
		}
	}
	return
}

func (model *Model) ApplicableProfiles() (applicableProfiles []data.Profile) {
	companyPoints := 0

	if model.Company != nil {
		companyPoints = model.Company.Points()
	}

	for _, profile := range model.Api.Cache.Profiles {
		if profile.Points+companyPoints <= model.Limits.CompanyMaximumPoints {
			applicableProfiles = append(applicableProfiles, profile)
		}
	}

	return
}

func (model *Model) RemovableTraits() (removableTraits []data.Trait) {
	traits := data.FilterTraitsByType("Special", model.Api.Cache.Traits)
	for _, trait := range traits {
		if utils.Contains(model.Traits, trait.Name) {
			removableTraits = append(removableTraits, trait)
		}
	}
	for _, trait := range model.TraitsWithChoices {
		removableTraits = append(removableTraits, *trait)
	}
	return
}

func (model *Model) CurrentTraitWithChoice() *data.Trait {
	return model.TraitsWithChoices[model.Indexes.CurrentTraitWithChoice]
}

func (model *Model) UpdateCurrentTraitWithChoiceName() {
	model.TraitsWithChoices[model.Indexes.CurrentTraitWithChoice] = model.TraitsWithChoices[model.Indexes.CurrentTraitWithChoice].TraitWithChoiceUpdatedName()
}

func (model *Model) CurrentTraitChoice() *data.TraitChoice {
	return model.CurrentTraitWithChoice().Choices[model.Indexes.CurrentChoice]
}

func (model *Model) UpdateCompanyWorkingPointTotal() {
	companyPoints := 0
	groupPoints := 0

	if model.Company != nil {
		companyPoints = model.Company.Points()
	}

	if model.Group != nil {
		groupPoints = model.Points
	}

	model.Temp.CompanyPoints = companyPoints + groupPoints
}
