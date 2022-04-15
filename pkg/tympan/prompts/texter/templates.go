package texter

import "github.com/erikgeiser/promptkit/textinput"

// The default template for text input; shows the prompt message followed by a blank line before the input field. If the
// input length is zero, it reports in bold orange text that the input is invalid because it cannot be empty. If any
// characters have been written, it instead displays that teh input is valid in bold blue text.
const TemplateDefault = `
{{- Bold .Prompt }}

{{ .Input -}}
{{- if not .Valid }} {{ Foreground "166" (Bold "Invalid: Cannot be empty.") }}
{{- else }} {{ Foreground "32" (Bold "Valid") }}
{{- end -}}
`

// This template is the same as the default except that the invalid input message is determined by the options passed to
// the instance.
const TemplateCustomInvalidMessage = `
{{- Bold .Prompt }} {{ .Input -}}
{{- if not .Valid }} {{ Foreground "166" (Bold InvalidMessage) }}
{{- else }} {{ Foreground "32" (Bold "Valid") }}
{{- end -}}
`

// The default result template returns the prompt followed by a space and the input highlighted in blue.
const DefaultResultTemplate = `
{{- print .Prompt " " (Foreground "32"  (Mask .FinalValue)) "\n" -}}
`

// This option overrides the default template the prompt uses to display to the terminal. For more information, see the
// docs for TextInput: https://pkg.go.dev/github.com/erikgeiser/promptkit/textinput#TextInput.Template
func WithTemplate(template string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.Template = template
	}
}

// This option overrides the default template the prompt uses to display results to the terminal. For more information,
// see the docs for TextInput:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/textinput#TextInput.ResultTemplate
func WithResultTemplate(template string) Option {
	return func(prompt *textinput.TextInput) {
		prompt.ResultTemplate = template
	}
}
