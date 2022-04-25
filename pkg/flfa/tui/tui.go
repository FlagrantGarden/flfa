package tui

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SharedModel struct {
	Api *flfa.Api
	compositor.Compositor
}

func (model *SharedModel) SaveConfig() tea.Cmd {
	err := model.Api.Tympan.SaveConfig()
	if err != nil {
		return model.RecordFatalError(err)
	}

	model.State = compositor.StateSavedConfiguration

	return nil
}

func Title(subtitle string, width int, options ...terminal.Option) string {
	settings := TerminalSettings(options...)
	barColor := lipgloss.Color("#948AE3")
	bar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		settings.AppliedExtraStyles("title").Foreground(barColor).Render("Flagrant Factions"),
		settings.AppliedExtraStyles("subtitle").Foreground(barColor).Render(subtitle),
	)
	coloredBar := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(bar)
	return coloredBar
}

type Sizes struct {
	Buffer  BoxSize
	Content BoxSize
	ViewBox BoxSize
}

type BoxSize struct {
	Height int
	Width  int
}

func (model *SharedModel) DisplaySizes(content string) (sizes Sizes) {
	sizes.Buffer.Height = model.Height
	sizes.Buffer.Width = model.Width
	sizes.Content.Height = lipgloss.Height(content)
	sizes.Content.Width = lipgloss.Width(content)
	maxHeight := 60
	maxWidth := 180
	sizes.ViewBox.Width = sizes.Buffer.Width - 2
	sizes.ViewBox.Height = sizes.Buffer.Height - 2
	if sizes.ViewBox.Width > maxWidth {
		sizes.ViewBox.Width = maxWidth
	}
	if sizes.ViewBox.Height > maxHeight {
		sizes.ViewBox.Height = maxHeight
	}
	if sizes.Content.Width > sizes.ViewBox.Width-2 {
		model.FatalError = fmt.Errorf(
			"Too wide; buffer is %d and content is %d",
			sizes.Buffer.Width, sizes.Content.Width)
	}
	if sizes.Content.Height > sizes.ViewBox.Height-2 {
		model.FatalError = fmt.Errorf(
			"Too tall; buffer is %d and content is %d",
			sizes.Buffer.Height, sizes.Content.Height)
	}
	return sizes
}

func (model *SharedModel) Display(subtitle string, body string, options ...terminal.Option) (view string) {
	if model.Width == 0 || model.Height == 0 {
		return "loading"
	}
	sizes := model.DisplaySizes(body)
	title := Title(subtitle, sizes.ViewBox.Width)
	body = lipgloss.Place(sizes.Content.Width+2, sizes.Content.Height+2, lipgloss.Left, lipgloss.Top, body)
	view = lipgloss.JoinVertical(lipgloss.Center, title, body)
	view = lipgloss.Place(sizes.ViewBox.Width, sizes.ViewBox.Height, lipgloss.Left, lipgloss.Top, view)
	view = model.TerminalSettings.RenderWithDynamicStyle("app_box", view)
	return lipgloss.Place(sizes.Buffer.Width, sizes.Buffer.Height, lipgloss.Center, lipgloss.Center, view)
}

func TerminalSettings(options ...terminal.Option) *terminal.Settings {
	settings := terminal.New(
		// Start from compositor to inherit required settings
		terminal.From(*compositor.DefaultTerminalSettings()),
		// ----- Colors ----- //
		terminal.WithBodyColor(lipgloss.Color("#F7F1FF")),
		terminal.WithLeadColor(lipgloss.Color("#5AD4E6")),
		terminal.WithSubtleColor(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#767676"}),
		terminal.WithExtraColor("application", lipgloss.Color("#948AE3")),     // Purple
		terminal.WithExtraColor("error", lipgloss.Color("#FC618D")),           // Pink error
		terminal.WithExtraColor("background", lipgloss.Color("#363537")),      // Very Dark Gray
		terminal.WithExtraColor("warning", lipgloss.Color("#FD9353")),         // Light Orange
		terminal.WithExtraColor("highlight", lipgloss.Color("#7BD88F")),       // Light Green
		terminal.WithExtraColor("adding", lipgloss.Color("#5AD4E6")),          // Light Blue
		terminal.WithExtraColor("removing", lipgloss.Color("#FD9353")),        // Light Orange
		terminal.WithExtraColor("selected_lead", lipgloss.Color("#5AD4E6")),   // Light Blue
		terminal.WithExtraColor("selected_body", lipgloss.Color("#FCE566")),   // Light Yellow
		terminal.WithExtraColor("company_description", lipgloss.Color("212")), // Pale Pink
		terminal.WithExtraColor("company_name", lipgloss.Color("99")),         // Royal Purple
		// -----Style Components ----- //
		// Shorthand equivalent to the HTML <strong> tag
		terminal.WithExtraStyle("strong", lipgloss.NewStyle().Bold(true)),
		// The Title text is bold with a hidden border so padding/margins take effect
		terminal.WithExtraStyle("title", lipgloss.NewStyle().Bold(true).Border(lipgloss.HiddenBorder(), true)),
		// The Subtitle text is italic with a hidden border so padding/margins take effect
		terminal.WithExtraStyle("subtitle", lipgloss.NewStyle().Italic(true).Border(lipgloss.HiddenBorder(), true)),
		// Help text is fainter than normal and has horizontal padding
		terminal.WithExtraStyle("help", lipgloss.NewStyle().Faint(true).Padding(0, 1)),
		// The appplication box has a double-lined border to contain everything else
		terminal.WithExtraStyle("app_box_border", lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), true)),
		// Blockquotes have a thick border on their left and are rendered in italics.
		terminal.WithExtraStyle(
			"blockquote",
			lipgloss.NewStyle().BorderLeft(true).MarginLeft(1).PaddingLeft(1).Italic(true).BorderStyle(lipgloss.ThickBorder()),
		),
		terminal.WithExtraStyle(
			"trait_description_paragraph",
			lipgloss.NewStyle().BorderLeft(true).PaddingLeft(1).MarginLeft(4).BorderStyle(lipgloss.NormalBorder()),
		),
		// Table headers always have a double border beneath them.
		terminal.WithExtraStyle("table_header", lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.DoubleBorder())),
		// Table entries always have a normal border beneath them.
		terminal.WithExtraStyle("table_entry", lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.NormalBorder())),
		// Groups in tables have very particular requirements to display neatly:
		terminal.WithExtraStyle("table_group_name", lipgloss.NewStyle().Width(20).Align(lipgloss.Left)),
		terminal.WithExtraStyle("table_group_profile_name", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_melee", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_missile", lipgloss.NewStyle().Width(15).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_move", lipgloss.NewStyle().Width(10).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_fighting_strength", lipgloss.NewStyle().Width(10).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_resolve", lipgloss.NewStyle().Width(5).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_toughness", lipgloss.NewStyle().Width(5).Align(lipgloss.Center)),
		terminal.WithExtraStyle("table_group_traits", lipgloss.NewStyle().Width(25).Align(lipgloss.Left)),
		// ----- Dynamic Styles ----- //
		terminal.WithDynamicStyle(
			"app_box", // The application box adds the border and colorizes it
			terminal.OverrideWithExtraStyle("app_box_border"),
			terminal.ColorizeBorderForeground("application"),
		),
		terminal.WithDynamicStyle(
			"app_name", // The application name should always be bold and colored per the application settings
			terminal.OverrideWithStyle(lipgloss.NewStyle().Bold(true)),
			terminal.ColorizeForeground("application"),
		),
		terminal.WithDynamicStyle(
			"company_name",
			terminal.OverrideWithExtraStyle("strong"),
			terminal.ColorizeForeground("company_name"),
		),
		terminal.WithDynamicStyle(
			"company_description",
			terminal.OverrideWithExtraStyle("blockquote"),
			terminal.ColorizeForeground("company_description"),
			terminal.ColorizeBorderForeground("company_description"),
		),
		terminal.WithDynamicStyle(
			"captain_trait_description",
			terminal.OverrideWithExtraStyle("trait_description_paragraph"),
			terminal.ColorizeForeground("highlight"),
		),
		terminal.WithDynamicStyle(
			"confirmation_emphasis", // Used in the message text to highlight a piece of info.
			terminal.OverrideWithExtraStyle("strong"),
			terminal.ColorizeForeground("highlight"),
		),
		terminal.WithDynamicStyle(
			"warning_emphasis", // Used to highlight text to emphasize dangerous or destructive actions
			terminal.OverrideWithExtraStyle("strong"),
			terminal.ColorizeForeground("warning"),
		),
		terminal.WithDynamicStyle(
			"selection_marking",
			terminal.OverrideWithStyle(lipgloss.NewStyle().Bold(true).Width(4)),
			terminal.ColorizeForeground("adding"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("removing"), "removing"),
		),
		terminal.WithDynamicStyle(
			"selection_marking_table",
			terminal.OverrideWithStyle(lipgloss.NewStyle().Bold(true).Width(4).BorderBottom(true).BorderStyle(lipgloss.HiddenBorder())),
			terminal.ColorizeForeground("adding"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("removing"), "removing"),
		),
		terminal.WithDynamicStyle(
			"add_selected_trait_name", // Used in prompts when selecting a trait to add
			terminal.OverrideWithExtraStyle("strong"),
			terminal.ColorizeForeground("adding"),
		),
		terminal.WithDynamicStyle(
			"remove_selected_trait_name", // Used in prompts when selecting a trait to remove
			terminal.OverrideWithExtraStyle("strong"),
			terminal.ColorizeForeground("removing"),
		),
		terminal.WithDynamicStyle(
			"unselected_trait_name", // Used in prompts when selecting a trait
			terminal.OverrideWithExtraStyle("strong"),
		),
		terminal.WithDynamicStyle(
			"selected_trait_effect", // Used in prompts when selecting a trait
			terminal.ColorizeForeground("selected_body"),
		),
		terminal.WithDynamicStyle(
			"unselected_trait_effect", // Used in prompts when selecting a trait
			terminal.OverrideWithStyle(lipgloss.NewStyle().Faint(true)),
		),
		terminal.WithDynamicStyle(
			"captain_group_name",
			terminal.OverrideWithExtraStyle("strong"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_name",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_name"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_profile_name",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_profile_name"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_melee",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_melee"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_missile",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_missile"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_move",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_move"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_fighting_strength",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_fighting_strength"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_resolve",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_resolve"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_toughness",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_toughness"),
		),
		terminal.WithDynamicStyle(
			"table_group_header_traits",
			terminal.OverrideWithExtraStyle("table_header"),
			terminal.OverrideWithExtraStyle("strong"),
			terminal.OverrideWithExtraStyle("table_group_traits"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_name",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_name"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "adding"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagsAreOn("selected", "removing"), "removing"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_profile_name",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_profile_name"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_melee",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_melee"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_missile",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_missile"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_move",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_move"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_fighting_strength",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_fighting_strength"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_resolve",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_resolve"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_toughness",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_toughness"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
		terminal.WithDynamicStyle(
			"table_group_entry_traits",
			terminal.OverrideWithExtraStyle("table_entry"),
			terminal.OverrideWithExtraStyle("table_group_traits"),
			terminal.ColorizeForegroundConditionally(terminal.IfFlagIsOn("selected"), "selected_body"),
		),
	)

	for _, option := range options {
		option(settings)
	}

	return settings
}

func SelectedAndRemovingFlags() map[string]terminal.Flag {
	flags := make(map[string]terminal.Flag)
	flags["selected"] = terminal.FlagOn
	flags["removing"] = terminal.FlagOn

	return flags
}
