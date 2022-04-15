package selector

import "github.com/erikgeiser/promptkit/selection"

// This option allows you to specify the choices to pass to the new selection prompt without having to convert them to
// Choice objects yourself.
func WithChoices(choices any) Option {
	return func(prompt *selection.Selection) {
		prompt.Choices = selection.Choices(choices)
	}
}

// This option allows you to append one or more choices to those already included in the prompt.
func WithAdditionalChoices(choices ...any) Option {
	return func(prompt *selection.Selection) {
		prompt.Choices = append(prompt.Choices, selection.Choices(choices)...)
	}
}
