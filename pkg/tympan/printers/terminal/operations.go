package terminal

import "github.com/charmbracelet/lipgloss"

// Operations are functions which operate on a style, returning the modified style.
// Used in the Apply* method of a Settings instance.
type Operation func(settings *Settings, style lipgloss.Style) lipgloss.Style

// A helper operation to inherit values from the specified style and appliy
// them to the dynamic style, overwriting existing definitions. Only values
// explicitly set on the style in argument will be applied.
//
// Margins, padding, and underlying string values are not inherited. To add them,
// use OverrideWithStyle/OverrideWithExtraStyle, the WithPadding* functions, or
// ensure they're in the base style.
//
// For example:
//     base := lipgloss.NewStyle().Bold(true)
//     new := lipgloss.NewStyle().Bold(false).Italic(true)
//     operation := InheritFromStyle(new)
//     updated := operation(base)
//     updated.Render("Example")
// The example text is bold and italic because the base style was set to bold and it
// inherited the italic style. Even though the style to inherit specified bold as false,
// the resulting style has it set to true.
func InheritFromStyle(styleToInherit lipgloss.Style) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		return style.Inherit(styleToInherit)
	}
}

// A helper operation to override with values from the specified style.
// Overriding inverts the logic of inheritance, extending the overriding style with
// any values from the base style which are not set on the overriding style.
//
// For example:
//     base := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32"))
//     new := lipgloss.NewStyle().Bold(false).Italic(true)
//     operation := InheritFromStyle(new)
//     updated := operation(base)
//     updated.Render("Example")
// The example text is italic and blue but not bold. Unlike with normal inheritance, the
// settings in the overriding style are preferred over the base style and any non-conflicting
// styles - in this case, color and italics - are merged from both.
func OverrideWithStyle(overridingStyle lipgloss.Style) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		return overridingStyle.Copy().Inherit(style)
	}
}

// A helper operation to inherit values from the specified extra style and appliy
// them to the dynamic style, overwriting existing definitions. Only values
// explicitly set on the style in argument will be applied.
//
// Margins, padding, and underlying string values are not inherited. To add them,
// use OverrideWithStyle/OverrideWithExtraStyle, the WithPadding* functions, or
// ensure they're in the base style.
func InheritFromExtraStyle(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		source, ok := settings.Styles.Extra[name]
		if ok {
			return style.Inherit(source)
		}
		return style
	}
}

// A helper operation to override with values from the specified extra style.
// Overriding inverts the logic of inheritance, extending the overriding style with
// any values from the base style which are not set on the overriding style.
func OverrideWithExtraStyle(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		source, ok := settings.Styles.Extra[name]
		if ok {
			return source.Copy().Inherit(style)
		}
		return style
	}
}

// A helper operation to set the text's background with the specified extra color.
func ColorizeBackground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.Background(color)
		}

		return style
	}
}

// A helper operation to set the text's foreground with the specified extra color.
func ColorizeForeground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.Foreground(color)
		}
		return style
	}
}

// A helper operation to set the border's background with the specified extra color(s).
// The arguments work as follows:
//
// With one argument, the argument is applied to all sides.
//
// With two arguments, the arguments are applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the arguments are applied to the top side, the
// horizontal sides, and the bottom side, in that order.
//
// With four arguments, the arguments are applied clockwise starting from the
// top side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments nothing will be set.
func ColorizeBorderBackground(names ...string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if len(names) < 1 || len(names) > 4 {
			return style
		}

		var colors []lipgloss.TerminalColor

		for _, name := range names {
			color, ok := settings.Colors.Extra[name]
			if !ok {
				return style
			}

			colors = append(colors, color)
		}

		if len(colors) != len(names) {
			return style
		}
		return style.BorderBackground(colors...)
	}
}

// A helper operation to set the border's foreground with the specified extra color(s).
// The arguments work as follows:
//
// With one argument, the argument is applied to all sides.
//
// With two arguments, the arguments are applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the arguments are applied to the top side, the
// horizontal sides, and the bottom side, in that order.
//
// With four arguments, the arguments are applied clockwise starting from the
// top side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments nothing will be set.
func ColorizeBorderForeground(names ...string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if len(names) < 1 || len(names) > 4 {
			return style
		}

		var colors []lipgloss.TerminalColor

		for _, name := range names {
			color, ok := settings.Colors.Extra[name]
			if !ok {
				return style
			}

			colors = append(colors, color)
		}

		if len(colors) != len(names) {
			return style
		}

		return style.BorderForeground(colors...)
	}
}

// A helper operation to set the border's bottom background with the specified extra color(s).
func ColorizeBorderBottomBackground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderBottomBackground(color)
		}

		return style
	}
}

// A helper operation to set the border's left background with the specified extra color(s).
func ColorizeBorderLeftBackground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderLeftBackground(color)
		}

		return style
	}
}

// A helper operation to set the border's right background with the specified extra color(s).
func ColorizeBorderRightBackground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderRightBackground(color)
		}

		return style
	}
}

// A helper operation to set the border's top background with the specified extra color(s).
func ColorizeBorderTopBackground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderTopBackground(color)
		}

		return style
	}
}

// A helper operation to set the border's bottom foreground with the specified extra color(s).
func ColorizeBorderBottomForeground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderBottomForeground(color)
		}

		return style
	}
}

// A helper operation to set the border's left foreground with the specified extra color(s).
func ColorizeBorderLeftForeground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderLeftForeground(color)
		}

		return style
	}
}

// A helper operation to set the border's right foreground with the specified extra color(s).
func ColorizeBorderRightForeground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderRightForeground(color)
		}

		return style
	}
}

// A helper operation to set the border's top foreground with the specified extra color(s).
func ColorizeBorderTopForeground(name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok {
			return style.BorderTopForeground(color)
		}

		return style
	}
}

// A helper operation to replace the base style of the dynamic style instead of building
// on the Primary style from settings; because operations are applied in the order they
// are specified, this should always be the first operation passed in a list of operations.
func SetBaseStyle(base lipgloss.Style) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		return base.Copy()
	}
}
