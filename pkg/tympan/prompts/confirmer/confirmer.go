package confirmer

import (
	"text/template"

	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/muesli/termenv"
)

// Options are functions which modify a Confirmation prompt. They provide a semantic way to both discover configuration
// options for a Confirmation prompt and to pass them dynamically as needed.
type Option func(prompt *confirmation.Confirmation)

// This option allows you to pass additional functions for the templates used in the prompt. Specify one or more
// with their name as the key in a map. This option will add the function if it is not already registered or overwrite
// an extended function if it already exists for the prompt. For more information, see the docs for Confirmation:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/confirmation#Confirmation.ExtendedTemplateFuncs
// and for template.FuncMap: https://pkg.go.dev/text/template#FuncMap
func WithExtendedTemplateFuncs(funcMap template.FuncMap) Option {
	return func(prompt *confirmation.Confirmation) {
		for k, v := range funcMap {
			prompt.ExtendedTemplateFuncs[k] = v
		}
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
//     keymap.Abort := []string{"esc"} // replace ctrl+c with escape to abort
//     confirmer.NewModel("What's your name?", confirmer.WithKeyMap(keymap))
func WithKeyMap(keymap confirmation.KeyMap) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.KeyMap = &keymap
	}
}

// This option allows you to override the default wrap mode for the prompt (promptkit.WordWrap). The default mode wraps
// the input at width, wrapping on last white space before the word which runs over the width so that words are not cut
// in the middle. The other built-in modes are HardWrap, which wraps at the specified width regardless of the text, and
// nil which disables wrapping. You can also supply your own wrap mode by specifying a function which takes an input
// string and width in and returns the wrapped string.
func WithWrapMode(mode promptkit.WrapMode) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.WrapMode = mode
	}
}

// This option allows you to override the default value of the prompt (unconfirmed) when instantiating the prompt.
func WithDefaultValue(defaultValue confirmation.Value) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.DefaultValue = defaultValue
	}
}

// This option allows you to override the default display words for confirming/disconfirming ("yes" and "no"
// respectively) to something else. This is a shorthand option which also replaces the default template and result
// template for you to ones which respect the newly chosen option texts. If you want to use different display words for
// the prompt _and_ a custom template, you will need to update your template accordingly and should not use this option.
func WithCustomAnswers(yesText string, noText string) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.ExtendedTemplateFuncs["Yes"] = func() string {
			return yesText
		}
		prompt.ExtendedTemplateFuncs["No"] = func() string {
			return noText
		}
		prompt.Template = TemplateCustomOptions
		prompt.ResultTemplate = ResultTemplateCustomOptions
	}
}

// This option allows you to override how colors are rendered. By default, the underlying prompt queries the terminal.
func WithColorProfile(profile termenv.Profile) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.ColorProfile = profile
	}
}

// This option inverts the colors of the template; orange for yes and blue for no. Use when confirming is the more
// dangerous or destructive option.
func WithInvertedColorTemplate() Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.Template = InvertedColorTemplate
		prompt.ResultTemplate = InvertedColorResultTemplate
	}
}

// Create a new confirmation prompt by specifying a message to ask the user to answer with yes or no and zero or more
// options to configure the prompt's behavior.
func New(message string, options ...Option) *confirmation.Confirmation {
	prompt := confirmation.New(message, confirmation.Undecided)

	prompt.Template = DefaultTemplate
	prompt.ResultTemplate = DefaultTemplate

	for _, option := range options {
		option(prompt)
	}

	return prompt
}

// This helper function immediately places the created confirmation prompt into a model, returning it.
func NewModel(message string, options ...Option) *confirmation.Model {
	return confirmation.NewModel(New(message, options...))
}
