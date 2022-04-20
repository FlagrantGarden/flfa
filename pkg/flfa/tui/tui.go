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
		settings.AppliedExtraStyles("Title").Foreground(barColor).Render("Flagrant Factions"),
		settings.AppliedExtraStyles("Subtitle").Foreground(barColor).Render(subtitle),
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
	view = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.Color("#948AE3")).Render(view)
	return lipgloss.Place(sizes.Buffer.Width, sizes.Buffer.Height, lipgloss.Center, lipgloss.Center, view)
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
