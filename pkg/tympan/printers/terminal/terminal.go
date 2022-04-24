package terminal

import (
	"encoding/json"

	"github.com/charmbracelet/lipgloss"
)

// Configuration for how a given object should be printed to the terminal
type Settings struct {
	// Provides a semantic way to retrieve lipgloss styles for use when printing to the terminal
	Styles Styles
	// Provides a semantic way to retrieve lipgloss colors for use when printing to the terminal
	Colors Colors
	// Provides a semantic way to retrieve and compile lipgloss styles from those included in these settings
	DynamicStyles DynamicStyleList
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
			Primary: lipgloss.NewStyle(),
			Extra:   make(map[string]lipgloss.Style),
		},
		Colors: Colors{
			Extra: make(map[string]lipgloss.TerminalColor),
		},
		DynamicStyles: make(map[string][]Operation),
		Flags:         make(map[string]Flag),
	}

	for _, option := range options {
		option(settings)
	}

	return settings
}

// Compile a lipgloss style from various options and return the result.
// By default, it copies the primary style and applies all options to
// that style in the order they are specified.
func (settings *Settings) Apply(operations ...Operation) lipgloss.Style {
	style := settings.Styles.Primary.Copy()

	for _, option := range operations {
		style = option(settings, style)
	}

	return style
}

// Compile a lipgloss style from various options and render text with that style.
func (settings *Settings) ApplyAndRender(text string, operations ...Operation) string {
	return settings.Apply(operations...).Render(text)
}

// Deep clones a given terminal settings object, returning a pointer to the clone.
func (settings *Settings) Copy() (*Settings, error) {
	marshalled, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	clone := &Settings{}

	if err = json.Unmarshal(marshalled, clone); err != nil {
		return nil, err
	}

	return clone, nil
}

// Return a dynamic style from the list of dynamic styles, applying the
// current value of the styles and colors stored as settings through a
// list of operations to generate a style.
//
// Because this method returns the style, you can modify or extend a
// defined dynamic style before using it.
func (settings *Settings) DynamicStyle(name string) lipgloss.Style {
	dynamicStyleOperations, ok := settings.DynamicStyles[name]
	if !ok {
		return lipgloss.NewStyle()
	}

	return settings.Apply(dynamicStyleOperations...)
}

// Compile a dynamic style from the list of dynamic styles and render
// a string directly.
func (settings *Settings) RenderWithDynamicStyle(name string, text string) string {
	return settings.DynamicStyle(name).Render(text)
}

// TODO: remember why I added this.
// func As(source Settings) Option {
// 	return func(settings *Settings) {
// 		clone := source
// 		settings = &clone
// 	}
// }
