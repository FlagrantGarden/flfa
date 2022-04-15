package data

import (
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/json"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/charmbracelet/lipgloss"
)

func (company *Company) ToJson(options ...json.Option) string {
	return json.StructToJson(company, options...)
}

func (company *Company) ToMarkdown() {}

func (company *Company) ToTerminal() {}

func CompanyTerminalSettings(options ...terminal.Option) terminal.Settings {
	companyNameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32"))
	captainNameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("29"))
	captainTraitNameStyle := lipgloss.NewStyle().Bold(true)
	captainTraitEffectStyle := lipgloss.NewStyle().
		BorderLeft(true).
		PaddingLeft(1).
		MarginLeft(4).
		Foreground(lipgloss.Color("11")).
		BorderStyle(lipgloss.ThickBorder()).
		Width(60)
	combinedOptions := []terminal.Option{
		terminal.WithPrimaryStyle(lipgloss.NewStyle()),
		terminal.WithExtraStyle("company_name", companyNameStyle),
		terminal.WithExtraStyle("captain_name", captainNameStyle),
		terminal.WithExtraStyle("captain_trait_name", captainTraitNameStyle),
		terminal.WithExtraStyle("captain_trait_effect", captainTraitEffectStyle),
	}

	combinedOptions = append(combinedOptions, options...)

	return *terminal.New(combinedOptions...)
}
