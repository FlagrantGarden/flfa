package data

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/json"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/charmbracelet/lipgloss"
)

func (trait *Trait) ToJson(options ...json.Option) string {
	return json.StructToJson(trait, options...)
}

func (trait *Trait) ToMarkdownTableEntry() {}

func (trait *Trait) MarkdownHeader() {}

func TraitMarkdownTable(traits ...Trait) {}

func TraitTerminalSettings(options ...terminal.Option) terminal.Settings {
	combinedOptions := []terminal.Option{
		terminal.WithPrimaryStyle(lipgloss.NewStyle()),
		terminal.WithExtraStyle("lead", lipgloss.NewStyle().Bold(true)),
		terminal.WithExtraStyle("body", lipgloss.NewStyle()),
		terminal.WithLeadColor(lipgloss.Color("32")),
		terminal.WithBodyColor(lipgloss.Color("11")),
	}

	combinedOptions = append(combinedOptions, options...)

	return *terminal.New(combinedOptions...)

}

func (trait *Trait) ToTerminalChoice(selected bool, leadWidth int, options ...terminal.Option) string {
	settings := TraitTerminalSettings(options...)
	leadStyle := settings.AppliedExtraStyles("lead")
	bodyStyle := settings.AppliedExtraStyles("body")

	var marking string
	if selected {
		leadStyle.Foreground(settings.Colors.Lead)
		bodyStyle.Foreground(settings.Colors.Body)
		marking = fmt.Sprintf("  %s ", leadStyle.Render("»"))
	} else {
		bodyStyle.Faint(true)
		marking = "    "
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		marking,
		leadStyle.Align(lipgloss.Left).Width(leadWidth).Render(trait.Name),
		trait.DisplayEffectBlock(bodyStyle, 80),
	)
}

func (trait Trait) DisplayEffectBlock(style lipgloss.Style, width int) string {
	var effectParagraphs []string
	for index, line := range strings.Split(trait.Effect, "\n\n") {
		if index == 0 {
			line = strings.Join(strings.Split(line, "\n"), " ")
			effectParagraphs = append(effectParagraphs, style.Copy().Width(width).Render(line))
		} else if len(line) != 0 {
			line = strings.Join(strings.Split(line, "\n"), " ")
			effectParagraphs = append(effectParagraphs, style.Copy().Width(width).Render(line))
		}
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		effectParagraphs...,
	)
}
