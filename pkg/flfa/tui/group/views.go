package group

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/charmbracelet/lipgloss"
)

func (model *Model) FormattedGroupName() string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32")).Render(model.Name)
}

func (model *Model) GroupEditingOverview() string {
	header := fmt.Sprintf("Editing the '%s' Group. Current Profile:", model.FormattedGroupName())
	table := data.DisplayGroupTerminal(model.Group.ToSlice())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		"",
		table,
	)
}
