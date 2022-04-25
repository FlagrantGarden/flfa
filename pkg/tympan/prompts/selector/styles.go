package selector

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
)

// A choice style function takes a selection choice as input and returns a string for displaying in the terminal. Any
// styling, coloring, or munging of a choice object before displaying must be done in a choice style function and passed
// to the prompt for the SelectedChoiceStyle or UnselectedChoiceStyle field.
type ChoiceStyleFunc func(choice *selection.Choice) string

// This option allows you to override the default choice style that is called to determine how to render the user's
// choice after they have selected one and ended the prompt.
func WithFinalChoiceStyle(styleFunction ChoiceStyleFunc) Option {
	return func(prompt *selection.Selection) {
		prompt.FinalChoiceStyle = styleFunction
	}
}

// This option allows you to override the default choice style that is called to determine how to render a choice when
// the user has tentatively selected it but not hit enter yet.
func WithSelectedChoiceStyle(styleFunction ChoiceStyleFunc) Option {
	return func(prompt *selection.Selection) {
		prompt.SelectedChoiceStyle = styleFunction
	}
}

// This option allows you to override the default choice style that is called to determine how to render all the choices
// a user is not tentatively selecting at that time.
func WithUnselectedChoiceStyle(styleFunction ChoiceStyleFunc) Option {
	return func(prompt *selection.Selection) {
		prompt.UnselectedChoiceStyle = styleFunction
	}
}

// This helper function returns a minimal default final choice style for the prompt, bolding and coloring blue the
// string representation of the final choice
func DefaultFinalChoiceStyle() ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("32")).
			Render(fmt.Sprintf("  %s", choice.String))
	}
}

// This helper function returns a minimal final choice style for the prompt, bolding and applying the specified coloring
// to the string representation of the final choice.
func ColorizedBasicFinalChoiceStyle(color lipgloss.TerminalColor) ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(color).
			Render(fmt.Sprintf("  %s", choice.String))
	}
}

// This helper function returns a minimal default selected choice style for the prompt, bolding and coloring blue the
// string representation of the selected choice and placing a "»" before the choice to more clearly identify that it is
// selected beyond only using color/font weight.
func DefaultSelectedChoiceStyle() ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("32")).
			Render(fmt.Sprintf("» %s", choice.String))
	}
}

// This helper function returns a minimal selected choice style for the prompt, bolding and applying the specified
// coloring to the string representation of the selected choice and placing a "»" before the choice to more clearly
// identify that it is selected beyond only using color/font weight.
func ColorizedBasicSelectedChoiceStyle(color lipgloss.TerminalColor) ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(color).
			Render(fmt.Sprintf("» %s", choice.String))
	}
}

// This helper function returns a minimal default unselected choice style for the prompt, writing the string
// representation of the unselected choice after two spaces to ensure alignment with selected choices is not off.
func DefaultUnselectedChoiceStyle() ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		return fmt.Sprintf("  %s", choice.String)
	}
}
