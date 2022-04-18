package tui

import (
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
	return func() tea.Msg {
		err := model.Api.Tympan.SaveConfig()
		if err != nil {
			return model.RecordFatalError(err)
		}
		model.State = compositor.StateSavedConfiguration
		return nil
	}
}

func Title(subtitle string, width int, options ...terminal.Option) string {
	settings := TerminalSettings(options...)
	barColor := lipgloss.Color("#948AE3")
	bar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		settings.AppliedExtraStyles("Title").Foreground(barColor).Render("Flagrant Factions"),
		settings.AppliedExtraStyles("Subtitle").Foreground(barColor).Render(subtitle),
	)
	coloredBar := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(bar)
	return coloredBar
}

func (model *SharedModel) Display(subtitle string, body string, options ...terminal.Option) (view string) {
	title := Title(subtitle, model.Width-2)
	body = lipgloss.Place(lipgloss.Width(body)+2, lipgloss.Height(body)+2, lipgloss.Left, lipgloss.Top, body)
	view = lipgloss.JoinVertical(lipgloss.Center, title, body)
	view = lipgloss.Place(180, 40, lipgloss.Left, lipgloss.Top, view)
	view = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.Color("#948AE3")).Render(view)
	return lipgloss.Place(model.Width, model.Height, lipgloss.Center, lipgloss.Center, view)
}

func TerminalSettings(options ...terminal.Option) *terminal.Settings {
	settings := &terminal.Settings{
		Styles: terminal.Styles{
			Primary: lipgloss.NewStyle(),
			Extra: map[string]lipgloss.Style{
				"Title":    lipgloss.NewStyle().Bold(true).Border(lipgloss.HiddenBorder(), true),
				"Subtitle": lipgloss.NewStyle().Italic(true).Border(lipgloss.HiddenBorder(), true),
				"Help":     lipgloss.NewStyle().Faint(true).Padding(0, 1),
			},
		},
		Colors: terminal.Colors{
			Body:   lipgloss.Color("#F7F1FF"),
			Lead:   lipgloss.Color("#5AD4E6"),
			Subtle: lipgloss.Color("#69676C"),
			Extra: map[string]lipgloss.TerminalColor{
				"Background":   lipgloss.Color("#363537"),
				"Error":        lipgloss.Color("#FC618D"),
				"Warning":      lipgloss.Color("#FD9353"),
				"Highlight":    lipgloss.Color("#7BD88F"),
				"SelectedLead": lipgloss.Color("#5AD4E6"),
				"SelectedBody": lipgloss.Color("#FCE566"),
			},
		},
	}

	for _, option := range options {
		option(settings)
	}

	return settings
}
