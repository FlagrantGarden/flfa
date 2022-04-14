package terminal

import "github.com/charmbracelet/lipgloss"

// Configuration for how a given object should be printed to the terminal
type Settings struct {
	// Provides a semantic way to retrieve lipgloss styles for use when printing to the terminal
	Styles Styles
	// Provides a semantic way to retrieve lipgloss colors for use when printing to the terminal
	Colors Colors
	// Provides a way to specify settings to change the behavior of terminal-rendering functions dynamically.
	Flags map[string]Flag
}

// An Option returns a function which modifies a Settings object. Options provide a friendlier UX for creating settings.
type Option func(settings *Settings)

// Creates a new Settings object, ensuring the maps for flags and extra colors/styles exist. Applies specified options
// in the order they are specified (if they conflict, last option applies).
func New(options ...Option) (settings *Settings) {
	settings = &Settings{
		Styles: Styles{
			Extra: make(map[string]lipgloss.Style),
		},
		Colors: Colors{
			Extra: make(map[string]lipgloss.TerminalColor),
		},
		Flags: make(map[string]Flag),
	}

	for _, option := range options {
		option(settings)
	}

	return settings
}

// TODO: remember why I added this.
// func As(source Settings) Option {
// 	return func(settings *Settings) {
// 		clone := source
// 		settings = &clone
// 	}
// }

// TODO: remember why I added this.
// Clones a given terminal settings object, returning a pointer to the new object.
// func (settings Settings) Copy() *Settings {
// 	copy := settings
// 	return &copy
// }
