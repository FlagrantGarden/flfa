package confirmer

import "github.com/erikgeiser/promptkit/confirmation"

// The default template for displaying the prompt and available choices. It shows the prompt message with a blank space
// before displaying yes and no. When yes is selected, it is highlighted in blue. When no is selected, it is highlighted
// in orange.
const DefaultTemplate = `
{{- Bold .Prompt }}

{{ if .YesSelected -}}
	{{- print "\t" (Foreground "32" (Bold " »Yes")) "\t  No" -}}
{{- else if .NoSelected -}}
	{{- print "\t  Yes\t" (Foreground "166" (Bold " »No")) -}}
{{- else -}}
	{{- "  Yes  No" -}}
{{- end -}}
`

// This template is the same as the default template except that the colors for yes and no are inverted; orange for yes
// and blue for no. Useful when you want to highlight that _confirming_ is the dangerous option.
const InvertedColorTemplate = `
{{- Bold .Prompt }}

{{ if .YesSelected -}}
	{{- print "\t" (Foreground "166" (Bold " »Yes")) "\t  No" -}}
{{- else if .NoSelected -}}
	{{- print "\t  Yes\t" (Foreground "32" (Bold " »No")) -}}
{{- else -}}
	{{- "  Yes  No" -}}
{{- end -}}
`

// The default template for displaying the results of the prompt. It shows the prompt followed by the choice after a
// space, highlighted in blue for yes and orange for no.
const DefaultResultTemplate = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- Foreground "32" "Yes" -}}
{{- else -}}
	{{- Foreground "166" "No" -}}
{{- end }}
`

// This template is the same as the default result template except that the colors for yes and no are inverted; orange
// for yes and blue for no. Useful when you want to highlight that _confirming_ is the dangerous option.
const InvertedColorResultTemplate = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- Foreground "166" "Yes" -}}
{{- else -}}
	{{- Foreground "32" "No" -}}
{{- end }}
`

// A helper template for when you want to use words other than "yes" and "no" to answer the prompt question. It is
// otherwise identical in behavior to the default template.
const TemplateCustomOptions = `
{{- Bold .Prompt -}}

{{ if .YesSelected -}}
	{{- print "\t" (Foreground "32" (Bold "  »" Yes )) "\t  " No -}}
{{- else if .NoSelected -}}
	{{- print "\t  " Yes "\t"(Foreground "166" (Bold " »" No)) -}}
{{- else -}}
	{{- print "\t  Yes\t  No" -}}
{{- end -}}
`

// A helper result template for when you want to use words other than "yes" and "no" to answer the prompt question. It
// is otherwise identical in behavior to the default result template.
const ResultTemplateCustomOptions = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- Foreground "32" Yes -}}
{{- else -}}
	{{- Foreground "166" No -}}
{{- end }}
`

// This option overrides the default template the prompt uses to display to the terminal. For more information, see the
// docs for Confirmation: https://pkg.go.dev/github.com/erikgeiser/promptkit/confirmation#Confirmation.Template
func WithTemplate(template string) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.Template = template
	}
}

// This option overrides the default template the prompt uses to display results to the terminal. For more information,
// see the docs for Confirmation:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/confirmation#Confirmation.ResultTemplate
func WithResultTemplate(template string) Option {
	return func(prompt *confirmation.Confirmation) {
		prompt.ResultTemplate = template
	}
}
