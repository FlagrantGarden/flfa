package prompts

import "github.com/charmbracelet/lipgloss"

func StyleLead(chosen bool) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(true)

	if chosen {
		return style.Foreground(lipgloss.Color("32"))
	}

	return style
}

func StyleBody(chosen bool) lipgloss.Style {
	style := lipgloss.NewStyle().Faint(true)

	if chosen {
		return style.Foreground(lipgloss.Color("11"))
	}

	return style
}
