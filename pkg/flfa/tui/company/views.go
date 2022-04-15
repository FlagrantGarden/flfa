package company

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/charmbracelet/lipgloss"
)

func (model *Model) CaptainSummary() string {
	group := model.CaptainsGroup()
	trait := group.Captain

	companyName := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("32")).
		Render(model.Company.Name)

	effect := lipgloss.NewStyle().
		BorderLeft(true).
		PaddingLeft(1).
		MarginLeft(4).
		Faint(true).
		Foreground(lipgloss.Color("11")).
		BorderStyle(lipgloss.NormalBorder()).
		Width(60).
		Render(strings.TrimSpace(trait.Effect))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("The %s Captain's Group is currently '%s'.", companyName, group.Name),
		fmt.Sprintf("They have the %s trait:", lipgloss.NewStyle().Bold(true).Render(trait.Name)),
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

	summary.WriteString(fmt.Sprintf("Editing the '%s' Company. Current Roster (%d points):\n\n", model.Name, model.Points()))
	summary.WriteString(data.DisplayGroupTerminal(model.Groups))
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
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Render(model.Name)
}

func (model *Model) FormattedCompanyDescription() string {
	return lipgloss.NewStyle().
		BorderLeft(true).
		MarginLeft(1).
		PaddingLeft(3).
		Width(80).
		Italic(true).
		Foreground(lipgloss.Color("212")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("212")).
		Render(model.Description)
}
