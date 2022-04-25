package group

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/charmbracelet/lipgloss"
)

func (model *Model) FormattedGroupName() string {
	return model.TerminalSettings.Apply(
		terminal.OverrideWithExtraStyle("strong"),
		terminal.ColorizeForeground("adding"),
	).Render(model.Name)
}

func (model *Model) GroupEditingOverview() string {
	header := fmt.Sprintf("Editing the '%s' Group. Current Profile:", model.FormattedGroupName())
	table := data.DisplayGroupTerminal(model.TerminalSettings, model.Group.ToSlice())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		"",
		table,
	)
}
