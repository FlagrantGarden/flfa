package selector

import (
	"text/template"

	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/muesli/termenv"
)

// Options are functions which modify a Selection prompt. They provide a semantic way to both discover configuration
// options for a Selection prompt and to pass them dynamically as needed.
type Option func(prompt *selection.Selection)

// A template name function is used to tell the prompt what the name of a choice is when choosing between structs.
type TemplateNameFunc func(choice *selection.Choice) string

// A template header row function is used to write the header row when using the table view for selection. It takes
// a boolean value for whether or not the user currently has the option to scroll up to prior choices in the table.
type TemplateHeaderRowFunc func(canScrollUp bool) string

// This option allows you to pass additional functions for the templates used in the prompt. Specify one or more
// with their name as the key in a map. This option will add the function if it is not already registered or overwrite
// an extended function if it already exists for the prompt. For more information, see the docs for Selection:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/selection#Selection.ExtendedTemplateFuncs
// and for template.FuncMap: https://pkg.go.dev/text/template#FuncMap
func WithExtendedTemplateFuncs(funcMap template.FuncMap) Option {
	return func(prompt *selection.Selection) {
		for k, v := range funcMap {
			prompt.ExtendedTemplateFuncs[k] = v
		}
	}
}

// This option allows you to specify how many choices to display at once. If the page size is smaller than the number of
// choices or is zero, pagination is disabled.
func WithPageSize(size int) Option {
	return func(prompt *selection.Selection) {
		prompt.PageSize = size
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
//     selector.New(
//       "How many apples do you have?",
//       "selection.Choices([]int{0, 1, 2, 3, 5}),
//       "selector.WithKeyMap(keymap),
//     )
func WithKeyMap(keymap selection.KeyMap) Option {
	return func(prompt *selection.Selection) {
		prompt.KeyMap = &keymap
	}
}

// This option allows you to override the default wrap mode for the prompt (promptkit.WordWrap). The default mode wraps
// the input at width, wrapping on last white space before the word which runs over the width so that words are not cut
// in the middle. The other built-in modes are HardWrap, which wraps at the specified width regardless of the text, and
// nil which disables wrapping. You can also supply your own wrap mode by specifying a function which takes an input
// string and width in and returns the wrapped string.
func WithWrapMode(mode promptkit.WrapMode) Option {
	return func(prompt *selection.Selection) {
		prompt.WrapMode = mode
	}
}

// This option allows you to override how colors are rendered. By default, the underlying prompt queries the terminal.
func WithColorProfile(profile termenv.Profile) Option {
	return func(prompt *selection.Selection) {
		prompt.ColorProfile = profile
	}
}

// Create a new selection prompt by specifying a message to help the user make a decision, a list of choices to choose
// from, and zero or more options to configure the prompt's behavior.
func New(message string, choices []*selection.Choice, options ...Option) *selection.Selection {
	prompt := selection.New(message, choices)

	prompt.Template = TemplateDefault
	prompt.ResultTemplate = ResultTemplateDefault
	prompt.SelectedChoiceStyle = DefaultSelectedChoiceStyle()
	prompt.UnselectedChoiceStyle = DefaultUnselectedChoiceStyle()
	prompt.FinalChoiceStyle = DefaultFinalChoiceStyle()

	for _, option := range options {
		option(prompt)
	}

	return prompt
}

// This helper function immediately places the created selection prompt into a model, returning it.
func NewModel(message string, choices []*selection.Choice, options ...Option) *selection.Model {
	return selection.NewModel(New(message, choices, options...))
}

// This helper function enables you to pass a slice of strings instead of needing to convert them into choices first.
func NewStringSelector(message string, choices []string, options ...Option) *selection.Selection {
	return New(message, selection.Choices(choices), options...)
}

// This helper function immediately places the created selection prompt generated from string choices into a model.
func NewStringModel(message string, choices []string, options ...Option) *selection.Model {
	return selection.NewModel(NewStringSelector(message, choices, options...))
}

// This helper function includes the minimum required options for a functional selection prompt where choices are
// structs. To use it, you need to pass the message and choices as normal. You also need to pass a filter function, to
// enable the prompt to filter for valid choices when the user types, and a name function, to tell the prompt what to
// display for each entry in the selection list.
func NewStructSelector(
	message string,
	choices any,
	filter FilterFunc,
	nameFunc TemplateNameFunc,
	options ...Option,
) *selection.Selection {
	var combinedOptions []Option
	combinedOptions = append(combinedOptions, WithResultTemplate(ResultTemplateByName))
	combinedOptions = append(combinedOptions, WithFilter(filter))
	combinedOptions = append(combinedOptions, WithExtendedTemplateFuncs(template.FuncMap{"Name": nameFunc}))
	combinedOptions = append(combinedOptions, options...)

	return New(
		message,
		selection.Choices(choices),
		combinedOptions...,
	)
}

// This helper function wraps the call to NewStructSelector and returns a model directly.
func NewStructModel(
	message string,
	choices any,
	filter FilterFunc,
	nameFunc TemplateNameFunc,
	options ...Option,
) *selection.Model {
	return selection.NewModel(NewStructSelector(message, choices, filter, nameFunc, options...))
}

// This helper function includes the minimum required options for a functional selection prompt displayed as a table. To
// use it, you need to pass the message and choices as normal. You also need to pass a filter function, name function,
// header row function, and choice style functions for when each row is selected/unselected. These choice style
// functions should return the table row for that choice, formatted appropriately.
func NewTableSelector(
	message string,
	choices any,
	filter FilterFunc,
	nameFunc TemplateNameFunc,
	headerRowFunc TemplateHeaderRowFunc,
	selectedChoiceStyle ChoiceStyleFunc,
	unselectedChoiceStyle ChoiceStyleFunc,
	options ...Option,
) *selection.Selection {
	var combinedOptions []Option
	combinedOptions = append(combinedOptions, WithTemplate(TemplateTable))
	combinedOptions = append(combinedOptions, WithExtendedTemplateFuncs(template.FuncMap{"HeaderRow": headerRowFunc}))
	combinedOptions = append(combinedOptions, WithSelectedChoiceStyle(selectedChoiceStyle))
	combinedOptions = append(combinedOptions, WithUnselectedChoiceStyle(unselectedChoiceStyle))
	combinedOptions = append(combinedOptions, options...)

	return NewStructSelector(message, choices, filter, nameFunc, combinedOptions...)
}

// This helper function wraps the call to NewTableSelector and returns a model directly.
func NewTableModel(
	message string,
	choices any,
	filter FilterFunc,
	nameFunc TemplateNameFunc,
	headerRowFunc TemplateHeaderRowFunc,
	selectedChoiceStyle ChoiceStyleFunc,
	unselectedChoiceStyle ChoiceStyleFunc,
	options ...Option,
) *selection.Model {
	return selection.NewModel(NewTableSelector(
		message,
		choices,
		filter,
		nameFunc,
		headerRowFunc,
		selectedChoiceStyle,
		unselectedChoiceStyle,
		options...,
	))
}
