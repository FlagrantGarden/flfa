package company

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/charmbracelet/lipgloss"
)

func (model *Model) CaptainSummary() string {
	group := model.CaptainsGroup()
	trait := group.Captain

	traitName := model.TerminalSettings.Apply(
		terminal.OverrideWithExtraStyle("strong"),
		terminal.ColorizeForeground("highlight"),
	).Render(trait.Name)

	groupName := model.TerminalSettings.Apply(
		terminal.OverrideWithExtraStyle("strong"),
		terminal.ColorizeForeground("adding"),
	).Render(group.Name)

	effect := model.TerminalSettings.DynamicStyle("captain_trait_description").Width(60).
		Render(strings.TrimSpace(trait.Effect))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("The %s Captain's Group is currently '%s'.", model.FormattedCompanyName(), groupName),
		fmt.Sprintf("They have the %s trait:", traitName),
		effect,
	)
}

func (model *Model) ShouldDisplayCaptainSummary() bool {
	requiringStates := []SubstateEditing{
		SelectingCaptainReplacement,
		SelectingCaptainOption,
		RerollingCaptainTrait,
		SelectingCaptainTrait,
		ConfirmingCaptainDemotion,
		ConfirmingCaptainReplacement,
	}

	for _, state := range requiringStates {
		if state == model.Substate.Editing {
			return true
		}
	}

	return false
}

func (model *Model) CompanyOverview() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Editing the %s Company.\n", model.FormattedCompanyName()))
	summary.WriteString(model.FormattedCompanyDescription())
	summary.WriteString(fmt.Sprintf("\nCurrent Roster (%d points):\n\n", model.Points()))
	summary.WriteString(data.DisplayGroupTerminal(model.TerminalSettings, model.Groups))
	summary.WriteString("\n\n")

	return summary.String()
}

func (model *Model) ShouldDisplayCompanyOverview() bool {
	requiringStates := []SubstateEditing{
		SelectingOption,
		Renaming,
		Redescribing,
	}

	for _, state := range requiringStates {
		if state == model.Substate.Editing {
			return true
		}
	}

	return false
}

func (model *Model) FormattedCompanyName() string {
	return model.TerminalSettings.RenderWithDynamicStyle("company_name", model.Name)
}

func (model *Model) FormattedCompanyDescription() string {
	return model.TerminalSettings.DynamicStyle("company_description").
		Width(80).Render(strings.TrimSpace(model.Description))
}
