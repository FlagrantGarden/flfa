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

func GroupTerminalSettings(options ...pterm.Option) (settings *pterm.Settings) {
	combinedOptions := []pterm.Option{
		pterm.WithPrimaryStyle(lipgloss.NewStyle().Padding(0, 1)),
		pterm.WithExtraStyle("selected_body", lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Faint(true)),
		pterm.WithExtraStyle("captain", lipgloss.NewStyle().Bold(true)),
		pterm.WithExtraStyle("header", lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.DoubleBorder())),
		pterm.WithExtraStyle("name", lipgloss.NewStyle().Width(20).Align(lipgloss.Left)),
		pterm.WithExtraStyle("profile_name", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		pterm.WithExtraStyle("melee", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		pterm.WithExtraStyle("missile", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		pterm.WithExtraStyle("move", lipgloss.NewStyle().Width(10).Align(lipgloss.Center)),
		pterm.WithExtraStyle("fighting_strength", lipgloss.NewStyle().Width(10).Align(lipgloss.Center)),
		pterm.WithExtraStyle("resolve", lipgloss.NewStyle().Width(5).Align(lipgloss.Center)),
		pterm.WithExtraStyle("toughness", lipgloss.NewStyle().Width(5).Align(lipgloss.Center)),
		pterm.WithExtraStyle("traits", lipgloss.NewStyle().Width(25).Align(lipgloss.Left)),
	}

	combinedOptions = append(combinedOptions, options...)

	settings = pterm.New(combinedOptions...)

	row := lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder())

	if settings.Colors.Subtle == nil {
		settings.Styles.Extra["row"] = row.BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#767676"})
	} else {
		settings.Styles.Extra["row"] = row.BorderForeground(settings.Colors.Subtle)
	}

	return settings
}

func (group *Group) ToTerminalTableEntry(options ...pterm.Option) string {
	settings := GroupTerminalSettings(options...)
	name := fmt.Sprintf("%s (%d)", group.Name, group.Points)

	var cells []string
	var body_styles []string
	var lead_styles []string

	if group.Captain.Name != "" {
		lead_styles = append(lead_styles, "captain")
	}

	switch settings.Flag("selected") {
	case pterm.FlagOff:
		cells = append(cells, settings.AppliedExtraStyles().Width(4).Render("    "))
	case pterm.FlagOn:
		lead_styles = append(lead_styles, "selected_lead")
		body_styles = append(body_styles, "selected_body")
		cells = append(cells, settings.AppliedExtraStyles("selected_lead").Width(4).Render("  » "))
	}

	cells = append(cells, settings.AppliedExtraStyles(append(lead_styles, "name")...).Render(name))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "profile_name")...).Render(group.ProfileName))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "melee")...).Render(group.Melee.String()))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "missile")...).Render(group.Missile.String()))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "move")...).Render(group.Move.String()))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "fighting_strength")...).Render(group.FightingStrength.String()))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "resolve")...).Render(fmt.Sprintf("%d", group.Resolve)))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "toughness")...).Render(fmt.Sprintf("%d", group.Toughness)))
	cells = append(cells, settings.AppliedExtraStyles(append(body_styles, "traits")...).Render(strings.Join(group.Traits, ", ")))

	entry := lipgloss.JoinHorizontal(lipgloss.Center, cells...)

	// make sure there's always space around the group
	height := lipgloss.Height(entry)
	entry = lipgloss.PlaceVertical((height + 2), lipgloss.Center, entry)

	return settings.Styles.Extra["row"].Render(entry)
}

func (group *Group) TableHeaderTerminal(options ...pterm.Option) string {
	settings := GroupTerminalSettings(options...)

	var cells []string

	if settings.Flag("for_selection") == pterm.FlagOn {
		if settings.Flag("can_scroll_up") == pterm.FlagOn {
			cells = append(cells, settings.AppliedExtraStyles("header").Width(4).Render(" ⇡  "))
		} else {
			cells = append(cells, settings.AppliedExtraStyles("header").Width(4).Render("    "))
		}
	}

	cells = append(cells, settings.AppliedExtraStyles("header", "name").Render("Name (P)"))
	cells = append(cells, settings.AppliedExtraStyles("header", "profile_name").Render("Profile"))
	cells = append(cells, settings.AppliedExtraStyles("header", "melee").Render("Melee"))
	cells = append(cells, settings.AppliedExtraStyles("header", "missile").Render("Missile"))
	cells = append(cells, settings.AppliedExtraStyles("header", "move").Render("Move"))
	cells = append(cells, settings.AppliedExtraStyles("header", "fighting_strength").Render("FS"))
	cells = append(cells, settings.AppliedExtraStyles("header", "resolve").Render("R"))
	cells = append(cells, settings.AppliedExtraStyles("header", "toughness").Render("T"))
	cells = append(cells, settings.AppliedExtraStyles("header", "traits").Render("Traits"))

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

func DisplayGroupTerminal(groups []Group) string {
	var rows []string

	header := groups[0].TableHeaderTerminal()
	rows = append(rows, header)

	for _, group := range groups {
		rows = append(rows, group.ToTerminalTableEntry())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)
}

func (group *Group) ToJson(options ...pjson.Option) string {
	return pjson.StructToJson(group, options...)
}
