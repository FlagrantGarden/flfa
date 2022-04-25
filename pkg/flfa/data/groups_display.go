package data

import (
	"fmt"
	"strings"

	pjson "github.com/FlagrantGarden/flfa/pkg/tympan/printers/json"
	pterm "github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/charmbracelet/lipgloss"
)

func (group *Group) ToMarkdownTableEntry() string {
	output := strings.Builder{}
	traits := group.Traits
	if group.Captain.Name == "" {
		output.WriteString(fmt.Sprintf("| %s |", group.Name))
	} else {
		output.WriteString(fmt.Sprintf("| **%s** |", group.Name))
		traits = append([]string{fmt.Sprintf("**%s**", group.Captain.Name)}, traits...)
	}
	output.WriteString(fmt.Sprintf(" %s |", group.ProfileName))
	output.WriteString(fmt.Sprintf(" %s |", group.Melee.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.Missile.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.Move.String()))
	output.WriteString(fmt.Sprintf(" %s |", group.FightingStrength.String()))
	output.WriteString(fmt.Sprintf(" %d+ |", group.Resolve))
	output.WriteString(fmt.Sprintf(" %d |", group.Toughness))
	output.WriteString(fmt.Sprintf(" %s |\n", strings.Join(traits, ", ")))
	return output.String()
}

func (group *Group) MarkdownHeader() string {
	header := strings.Builder{}
	header.WriteString("| Name | Profile | Melee | Missile | Move | FS | R | T | Traits |\n")
	header.WriteString("| ---- | ------- | ----- | ------- | ---- | -- | - | - | ------ |\n")

	return header.String()
}

func GroupMarkdownTable(groups ...Group) string {
	groupOutput := strings.Builder{}
	groupOutput.WriteString("| Name | Profile | Melee | Missile | Move | FS | R | T | Traits |\n")
	groupOutput.WriteString("| ---- | ------- | ----- | ------- | ---- | -- | - | - | ------ |\n")
	for _, group := range groups {
		groupOutput.WriteString(group.ToMarkdownTableEntry())
	}
	return groupOutput.String()
}

func (group *Group) ToTerminalTableEntry(settings *pterm.Settings) string {
	name := fmt.Sprintf("%s (%d)", group.Name, group.Points)

	var cells []string

	if settings.FlagIsOn("selected") {
		cells = append(cells, settings.RenderWithDynamicStyle("selection_marking_table", "  » "))
	} else if settings.FlagIsOff("selected") {
		cells = append(cells, settings.RenderWithDynamicStyle("selection_marking_table", "    "))
	}

	nameStyle := settings.DynamicStyle("table_group_entry_name")
	if group.Captain.Name != "" {
		nameStyle.Bold(true)
	}
	cells = append(cells, nameStyle.Render(name))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_profile_name", group.ProfileName))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_melee", group.Melee.String()))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_missile", group.Missile.String()))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_move", group.Move.String()))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_fighting_strength", group.FightingStrength.String()))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_resolve", fmt.Sprintf("%d", group.Resolve)))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_toughness", fmt.Sprintf("%d", group.Toughness)))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_entry_traits", strings.Join(group.Traits, ", ")))

	entry := lipgloss.JoinHorizontal(lipgloss.Center, cells...)

	// make sure there's always space around the group
	height := lipgloss.Height(entry)
	if height%2 == 0 {
		height++
	}
	entry = lipgloss.PlaceVertical((height), lipgloss.Center, entry)

	return entry

	// return settings.Styles.Extra["row"].Render(entry)
}

func (group *Group) TableHeaderTerminal(settings *pterm.Settings) string {
	var cells []string

	if settings.Flag("for_selection") == pterm.FlagOn {
		scrollerStyle := lipgloss.NewStyle().Width(4).BorderBottom(true).BorderStyle(lipgloss.HiddenBorder())
		if settings.Flag("can_scroll_up") == pterm.FlagOn {
			cells = append(cells, scrollerStyle.Render("⇡   "))
		} else {
			cells = append(cells, scrollerStyle.Render("    "))
		}
	}

	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_name", "Name (P)"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_profile_name", "Profile"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_melee", "Melee"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_missile", "Missile"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_move", "Move"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_fighting_strength", "FS"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_resolve", "R"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_toughness", "T"))
	cells = append(cells, settings.RenderWithDynamicStyle("table_group_header_traits", "Traits"))

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

func DisplayGroupTerminal(settings *pterm.Settings, groups []Group) string {
	var rows []string

	header := groups[0].TableHeaderTerminal(settings)
	rows = append(rows, header)

	for _, group := range groups {
		rows = append(rows, group.ToTerminalTableEntry(settings))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)
}

func (group *Group) ToJson(options ...pjson.Option) string {
	return pjson.StructToJson(group, options...)
}
