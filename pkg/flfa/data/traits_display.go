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

func (trait *Trait) ToTerminalChoice(settings *terminal.Settings, leadWidth int) string {
	var marking string
	var nameStyle lipgloss.Style
	var effectStyle lipgloss.Style

	selected := settings.FlagIsOn("selected")
	removing := settings.FlagIsOn("removing_trait")

	if selected {
		effectStyle = settings.DynamicStyle("selected_trait_effect")
		if removing {
			nameStyle = settings.DynamicStyle("remove_selected_trait_name")
		} else {
			nameStyle = settings.DynamicStyle("add_selected_trait_name")
		}
		marking = fmt.Sprintf("  %s ", nameStyle.Render("»"))
	} else {
		effectStyle = settings.DynamicStyle("unselected_trait_effect")
		nameStyle = settings.DynamicStyle("unselected_trait_name")
		marking = "    "
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		marking,
		nameStyle.Align(lipgloss.Left).Width(leadWidth).Render(trait.Name),
		trait.DisplayEffectBlock(effectStyle, 80),
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
