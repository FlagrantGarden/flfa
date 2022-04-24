package terminal

import "github.com/charmbracelet/lipgloss"

// The Colors struct provides a semantic way to access various colors for terminal settings.
type Colors struct {
	// Use the Subtle color whenever you want to deemphasize text.
	Subtle lipgloss.TerminalColor
	// Use the Lead color for headings and other more-important text.
	Lead lipgloss.TerminalColor
	// Use the Body color as the default for general text.
	Body lipgloss.TerminalColor
	// Use extra to hold an custom list of colors by name for use.
	Extra map[string]lipgloss.TerminalColor
}

// This terminal settings option sets the subtle color for a new instance of terminal settings.
func WithSubtleColor(color lipgloss.TerminalColor) Option {
	return func(settings *Settings) *Settings {
		settings.Colors.Subtle = color
		return settings
	}
}

// This terminal settings option sets the lead color for a new instance of terminal settings.
func WithLeadColor(color lipgloss.TerminalColor) Option {
	return func(settings *Settings) *Settings {
		settings.Colors.Lead = color
		return settings
	}
}

// This terminal settings option sets the body color for a new instance of terminal settings.
func WithBodyColor(color lipgloss.TerminalColor) Option {
	return func(settings *Settings) *Settings {
		settings.Colors.Body = color
		return settings
	}
}

// This terminal settings option adds an extra color to a new instance of terminal settings or updates the color if it
// has already been added to the map. This can be used to extend/modify the list of extra colors safely in a loop.
func WithExtraColor(name string, color lipgloss.TerminalColor) Option {
	return func(settings *Settings) *Settings {
		settings.Colors.Extra[name] = color
		return settings
	}
}

// This terminal settings option sets the extra colors for a new instance of terminal settings to exactly the map that
// you pass to it; if this option is applied after any others which modify the extra colors, it replaces all of them.
func WithExtraColors(extras map[string]lipgloss.TerminalColor) Option {
	return func(settings *Settings) *Settings {
		settings.Colors.Extra = extras
		return settings
	}
}

// Return an extra color from an instance of terminal settings by name.
func (settings *Settings) ExtraColor(name string) lipgloss.TerminalColor {
	return settings.Colors.Extra[name]
}
