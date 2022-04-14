package terminal

import "github.com/charmbracelet/lipgloss"

// The Colors struct provides a semantic way to access various colors for terminal settings.
type Styles struct {
	// Use the primary style as your default
	Primary lipgloss.Style
	// Use extra styles to hold a named set of styles to apply
	Extra map[string]lipgloss.Style
}

// This terminal settings option sets the primary style for a new instance of terminal settings.
func WithPrimaryStyle(style lipgloss.Style) Option {
	return func(settings *Settings) {
		settings.Styles.Primary = style
	}
}

// This terminal settings option adds an extra style to a new instance of terminal settings or updates the style if it
// has already been added to the map. This can be used to extend/modify the list of extra styles safely in a loop.
func WithExtraStyle(name string, style lipgloss.Style) Option {
	return func(settings *Settings) {
		settings.Styles.Extra[name] = style
	}
}

// This terminal settings option sets the extra styles for a new instance of terminal settings to exactly the map that
// you pass to it; if this option is applied after any others which modify the extra styles, it replaces all of them.
func WithExtraStyles(extras map[string]lipgloss.Style) Option {
	return func(settings *Settings) {
		settings.Styles.Extra = extras
	}
}

// This terminal settings option takes multiple styles and merges them into one, adding them as an extra style with the
// given name to a new instance of terminal settings. If the name already exists in the list of extra styles, the new
// styles will be inherited to it, extending the style but not overwriting it.
func WithMergedExtraStyles(name string, styles ...lipgloss.Style) Option {
	return func(settings *Settings) {
		base, ok := settings.Styles.Extra[name]
		if ok {
			settings.Styles.Extra[name] = MergeStyles(base, styles...)
		} else {
			settings.Styles.Extra[name] = CombineStyles(styles[0], styles[1:]...)
		}
	}
}

// This helper function copies the base style and inherits the remaining styles in order, extending the copied style
// but not overwriting any settings, before returning the combined style. The additional styles are inherited in the
// order they are given, so the first style to specify a setting (base or additional) wins out. Because the styles are
// applied to a copy of the base style, they do not modify it at all.
func CombineStyles(base lipgloss.Style, additional ...lipgloss.Style) (style lipgloss.Style) {
	style = base.Copy()

	for _, inheriting := range additional {
		style.Inherit(inheriting)
	}

	return style
}

// This helper function starts with the base style and inherits the remaining styles in order, extending the base style
// but not overwriting any settings, before returning the merged style. The additional styles are inherited in the
// order they are given, so the first style to specify a setting (base or additional) wins out. Because the styles are
// applied to the base style, they modify it in place as well as return the merged style to the caller.
func MergeStyles(base lipgloss.Style, additional ...lipgloss.Style) lipgloss.Style {
	for _, inheriting := range additional {
		base.Inherit(inheriting)
	}

	return base
}

// The AppliedExtraStyles method copies the primary style and then merges all of the specified extra styles from the
// map of extra styles in the Settings, creating and returning a new style made from inheriting the extra styles in the
// order they were specified onto a copy of the Primary style.
func (settings *Settings) AppliedExtraStyles(extras ...string) (appliedStyle lipgloss.Style) {
	appliedStyle = settings.Styles.Primary.Copy()

	for _, extra := range extras {
		if style, ok := settings.Styles.Extra[extra]; ok {
			appliedStyle = MergeStyles(appliedStyle, style)
		}
	}

	return appliedStyle
}
