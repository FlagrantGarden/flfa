package terminal

import "github.com/charmbracelet/lipgloss"

// Operations are functions which operate on a style, returning the modified style.
// Used in the Apply* method of a Settings instance.
type Operation func(settings *Settings, style lipgloss.Style) lipgloss.Style

// Conditions are functions which use the settings to determine whether or not an
// Operation should be performed on a dynamic style.
type Condition func(settings *Settings) bool

// Helper condition to check if a flag has a particular setting before applying an operation.
// Returns false if the flag does not have the specified setting.
func IfFlagIs(name string, setting Flag) Condition {
	return func(settings *Settings) bool {
		return settings.Flag(name) == setting
	}
}

// Helper condition to check if a list of flags have the specified settings before applying an operation.
// Loops over each flag in the map, returning false if any flag has the wrong setting.
func IfFlagsAre(flags map[string]Flag) Condition {
	return func(settings *Settings) bool {
		for name, setting := range flags {
			if settings.Flag(name) != setting {
				return false
			}
		}

		return true
	}
}

// Helper condition to check if a flag is set to on before applying an operation.
func IfFlagIsOn(name string) Condition {
	return func(settings *Settings) bool {
		return settings.FlagIsOn(name)
	}
}

// Helper condition to check if all of the specified flags are set to on before applying an operation.
// Returns false if any flag specified is set to off or is unset.
func IfFlagsAreOn(names ...string) Condition {
	return func(settings *Settings) bool {
		for _, name := range names {
			if !settings.FlagIsOn(name) {
				return false
			}
		}
		return true
	}
}

// Helper condition to check if a flag is set to off before applying an operation.
func IfFlagsIsOff(name string) Condition {
	return func(settings *Settings) bool {
		return settings.FlagIsOff(name)
	}
}

// Helper condition to check if all of the specified flags are set to off before applying an operation.
// Returns false if any flag specified is set to on or is unset.
func IfFlagsAreOff(names ...string) Condition {
	return func(settings *Settings) bool {
		for _, name := range names {
			if !settings.FlagIsOff(name) {
				return false
			}
		}
		return true
	}
}

// Helper condition to check if a flag is unset before applying an operation.
func IfFlagsIsUnset(name string) Condition {
	return func(settings *Settings) bool {
		return settings.FlagIsUnset(name)
	}
}

// Helper condition to check if all of the specified flags are unset before applying an operation.
// Returns false if any flag specified is set to on or off.
func IfFlagsAreUnset(names ...string) Condition {
	return func(settings *Settings) bool {
		for _, name := range names {
			if !settings.FlagIsUnset(name) {
				return false
			}
		}
		return true
	}
}

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

// As InheritFromStyle, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func InheritFromStyleConditionally(condition Condition, styleToInherit lipgloss.Style) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if condition(settings) {
			return style.Inherit(styleToInherit)
		}

		return style
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

// As OverrideWithStyle, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func OverrideWithStyleConditionally(condition Condition, overridingStyle lipgloss.Style) Operation {
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

// As InheritFromExtraStyle, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func InheritFromExtraStyleConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		source, ok := settings.Styles.Extra[name]
		if ok && condition(settings) {
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

// As OverrideWithExtraStyle, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func OverrideWithExtraStyleConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		source, ok := settings.Styles.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBackgroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeForegroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderBackgroundConditionally(condition Condition, names ...string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if !condition(settings) {
			return style
		}

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

// As ColorizeBorderForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderForegroundConditionally(condition Condition, names ...string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if !condition(settings) {
			return style
		}

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

// As ColorizeBorderBottomBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderBottomBackgroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderLeftBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderLeftBackgroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderRightBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderRightBackgroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderTopBackground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderTopBackgroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderBottomForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderBottomForegroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderLeftForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderLeftForegroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderRightForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderRightForegroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As ColorizeBorderTopForeground, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func ColorizeBorderTopForegroundConditionally(condition Condition, name string) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		color, ok := settings.Colors.Extra[name]
		if ok && condition(settings) {
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

// As SetBaseStyle, but only applied if the specified condition returns true.
// Conditions are checked against the settings when using Apply or calling a DynamicStyle;
// they don't need to be true when the condition is specified, only when it is checked.
func SetBaseStyleConditionally(condition Condition, base lipgloss.Style) Operation {
	return func(settings *Settings, style lipgloss.Style) lipgloss.Style {
		if condition(settings) {
			return base.Copy()
		}

		return style
	}
}
