package texter

import (
	"html/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/muesli/termenv"
)

// Options are functions which modify a TextInput prompt. They provide a semantic way to both discover configuration
// options for a TextInput prompt and to pass them dynamically as needed.
type Option func(prompt *textinput.TextInput)

// This option overrides the default placeholder text the prompt displays when the user has not typed any text yet. For
// more information, see the docs for TextInput:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/textinput#TextInput.Placeholder
func WithPlaceholder(placeholder string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.Placeholder = placeholder
	}
}

// This option sets the initial value to be passed to the input, enabling them to hit enter to accept or to edit the
// value first. For more information, see the docs for TextInput:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/textinput#TextInput.InitialValue
func WithInitialValue(value string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InitialValue = value
	}
}

// This option overrides the default validation check for the prompt (the default only checks that the input is not
// empty), preventing submission if the input does not pass validation. If you pass nil for this option, no validation
// is performed. To use this option, pass a function which takes an input string and returns true if it is valid and
// false if it is not.
func WithValidateFunc(validateFunc func(string) bool) Option {
	return func(prompt *textinput.TextInput) {
		prompt.Validate = validateFunc
	}
}

// This option is a helper option for replacing the message from the default validation check with one specific to the
// validation actually being performed. Only use this if you are also changing the default validation for the prompt. Do
// not use this and WithValidationMessageFunc() together as they are different implementations of the same thing.
func WithValidationMessage(message string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.ExtendedTemplateFuncs["InvalidMessage"] = func() string {
			return message
		}
	}
}

// This option is a helper option for replacing the message from the default validation check with a dynamic message
// that can be aware of why the validation failed. Only use this if you are also changing the default validation for the
// prompt. Do not use this and WithValidationMessage() together as they are different implementations of the same thing.
func WithValidationMessageFunc(messageFunc func() string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.ExtendedTemplateFuncs["InvalidMessage"] = messageFunc
	}
}

// This option overrides whether the prompt should treat the input as secret or not. If set to true, the prompt will
// mask the input to prevent anyone from reading it off the screen.
func WithHidden(hidden bool) Option {
	return func(prompt *textinput.TextInput) {
		prompt.Hidden = hidden
	}
}

// This option overrides the default hide mask for the prompt's input if it is marked as hidden. The mask will be
// rendered to the screen once for every rune the user inputs.
func WithHideMask(mask rune) Option {
	return func(prompt *textinput.TextInput) {
		prompt.HideMask = mask
	}
}

// This option sets a maximum length to the input for the prompt. By default, the prompt does not have a limit.
func WithCharLimit(limit int) Option {
	return func(prompt *textinput.TextInput) {
		prompt.CharLimit = limit
	}
}

// This option defines the maximum number of characters that the input can display at a time. If a user types in more
// characters than the width, the viewport for the input prompt will scroll with their cursor.
func WithInputWidth(width int) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InputWidth = width
	}
}

// This option allows you to pass additional functions for the templates used in the prompt. Specify one or more
// with their name as the key in a map. This option will add the function if it is not already registered or overwrite
// an extended function if it already exists for the prompt. For more information, see the docs for TextInput:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/textinput#TextInput.ExtendedTemplateFuncs
// and for template.FuncMap: https://pkg.go.dev/text/template#FuncMap
func WithExtendedTemplateFuncs(funcMap template.FuncMap) Option {
	return func(prompt *textinput.TextInput) {
		for k, v := range funcMap {
			prompt.ExtendedTemplateFuncs[k] = v
		}
	}
}

// This option allows you to override the default style for the input text itself. Pass a lipgloss style and it will be
// applied inline to to the input text.
func WithInputTextStyle(style lipgloss.Style) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InputTextStyle = style
	}
}

// This option allows you to override the default style for the background of the input text box itself. Pass a lipgloss
// style and it will be applied inline to to the input text box.
func WithInputBackgroundStyle(style lipgloss.Style) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InputBackgroundStyle = style
	}
}

// This option allows you to override the default style for the plcaseholder text itself. Pass a lipgloss style and it
// will be applied inline to to the placeholder text.
func WithInputPlaceholderStyle(style lipgloss.Style) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InputPlaceholderStyle = style
	}
}

// This option allows you to override the default style for the cursor when inputting text to the prompt. Pass a
// lipgloss style and it will be applied inline to to the cursor.
func WithInputCursorStyle(style lipgloss.Style) Option {
	return func(prompt *textinput.TextInput) {
		prompt.InputCursorStyle = style
	}
}

// This option allows you to override the default key map for the prompt. Specify the key presses (or combinations) you
// want to trigger an action from as strings for their specified action. Note that this option entirely replaces the
// existing key map; it is best practice to create the default key map, modify it, and then pass the modified map to
// this option instead of creating the key map inline.
//
// For example:
//
//     keymap := confirmation.NewDefaultKeyMap() // start with default map
//     keymap.Abort := []string{"shift+esc"} // replace ctrl+c with shift+escape to abort
//     texter.NewModel("Do you want to continue?", texter.WithKeyMap(keymap))
func WithKeyMap(keymap textinput.KeyMap) Option {
	return func(prompt *textinput.TextInput) {
		prompt.KeyMap = &keymap
	}
}

// This option allows you to override the default wrap mode for the prompt (promptkit.WordWrap). The default mode wraps
// the input at width, wrapping on last white space before the word which runs over the width so that words are not cut
// in the middle. The other built-in modes are HardWrap, which wraps at the specified width regardless of the text, and
// nil which disables wrapping. You can also supply your own wrap mode by specifying a function which takes an input
// string and width in and returns the wrapped string.
func WithWrapMode(mode promptkit.WrapMode) Option {
	return func(prompt *textinput.TextInput) {
		prompt.WrapMode = mode
	}
}

// This option allows you to override how colors are rendered. By default, the underlying prompt queries the terminal.
func WithColorProfile(profile termenv.Profile) Option {
	return func(prompt *textinput.TextInput) {
		prompt.ColorProfile = profile
	}
}

// Create a new textinput prompt by specifying a message to display to the user to explain what text they should input.
// Includes a default placeholder and validation message reminding them that the input cannot be empty. Any options
// passed to this function are applied in the order they are specified and after the defaults are set.
func New(message string, options ...Option) *textinput.TextInput {
	prompt := textinput.New(message)

	prompt.Template = TemplateDefault
	prompt.ResultTemplate = DefaultResultTemplate

	for _, option := range options {
		option(prompt)
	}

	return prompt
}

// This helper function provides a shorthand for creating a textinput prompt with a custom validation function. It is
// functionally identical to the New() function except that it requires a validation function (leading to slightly
// improved UX when developing and using a language server) and prepends the validation function to the list of any
// further options.
func NewValidatable(message string, validateFunc func(s string) bool, options ...Option) *textinput.TextInput {
	var combinedOptions []Option
	combinedOptions = append(combinedOptions, WithValidateFunc(validateFunc))
	combinedOptions = append(combinedOptions, options...)

	return New(
		message,
		combinedOptions...,
	)
}

// This helper function provides a shorthand for creating a textinput prompt with a custom validation function. It is
// functionally identical to the New() function except that it requires a validation function and string for the message
// to display when the input fails validation (leading to slightly improved UX when developing and using a language
// server) and prepends the validation function to the list of any further options.
func NewValidatableWithCustomMessage(
	message string,
	validateFunc func(s string) bool,
	validateMessage string,
	options ...Option,
) *textinput.TextInput {
	var combinedOptions []Option
	combinedOptions = append(combinedOptions, WithValidateFunc(validateFunc))
	combinedOptions = append(combinedOptions, WithValidationMessage(validateMessage))
	combinedOptions = append(combinedOptions, options...)

	return New(
		message,
		combinedOptions...,
	)
}

// This helper function provides a shorthand for creating a textinput prompt with a custom validation function. It is
// functionally identical to the New() function except that it requires a validation function and function for
// determiningthe message to display when the input fails validation (leading to slightly improved UX when developing
// and using a language server) and prepends the validation function to the list of any further options.
func NewValidatableWithCustomMessageFunc(
	message string,
	validateFunc func(s string) bool,
	validateMessageFunc func() string,
	options ...Option,
) *textinput.TextInput {
	var combinedOptions []Option
	combinedOptions = append(combinedOptions, WithValidateFunc(validateFunc))
	combinedOptions = append(combinedOptions, WithValidationMessageFunc(validateMessageFunc))
	combinedOptions = append(combinedOptions, options...)

	return New(
		message,
		combinedOptions...,
	)
}
