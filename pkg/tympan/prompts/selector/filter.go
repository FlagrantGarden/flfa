package selector

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
)

// A Filter function is a function which takes a filter string and a choice object as input and returns true if the
// choice matches the filter or false if it does not.
type FilterFunc func(filter string, choice *selection.Choice) bool

// This option enables you to pass your own custom filter function to the prompt, overriding the default, which compares
// the filter string to the string representation of the choice.
func WithFilter(filter FilterFunc) Option {
	return func(prompt *selection.Selection) {
		prompt.Filter = filter
	}
}

// This option enables you to override the default filter prompt, which reads "Filter:"
func WithFilterPrompt(message string) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterPrompt = message
	}
}

// This option enables you to override the default filter placeholder, which reads "Type to filter choices."
func WithFilterPlaceholder(placeholder string) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterPlaceholder = placeholder
	}
}

// This option allows you to override the default style for the filter's input text. Pass a lipgloss style and it will
// be applied inline to to the input text.
func WithFilterInputTextStyle(style lipgloss.Style) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterInputTextStyle = style
	}
}

// This option allows you to override the default style for the background of the filter's input text. Pass a lipgloss
// style and it will be applied inline to to the background.
func WithFilterInputBackgroundStyle(style lipgloss.Style) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterInputBackgroundStyle = style
	}
}

// This option allows you to override the default style for the filter's placeholder text. Pass a lipgloss style and it
// will be applied inline to to the input text.
func WithFilterInputPlaceholderStyle(style lipgloss.Style) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterInputPlaceholderStyle = style
	}
}

// This option allows you to override the default style for the cursor when inputting text for the prompt's filter. Pass
// a lipgloss style and it will be applied inline to to the cursor.
func WithFilterInputCursorStyle(style lipgloss.Style) Option {
	return func(prompt *selection.Selection) {
		prompt.FilterInputCursorStyle = style
	}
}
