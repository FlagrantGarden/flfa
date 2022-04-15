package selector

import "github.com/erikgeiser/promptkit/selection"

// The default template for selection; shows the bolded prompt message followed on the next line by the filter prompt
// (if the prompt is using a filter) and then a blank line before the selection list. If the list is paginated and the
// user can scroll up, a "⇡" is inserted immediately before the list. If the list is paginated and the user can scroll
// down, a "⇣" is inserted after the last visible choice in the list. As a user moves through the list, the tentatively
// selected choice becomes bolded and highlighted in blue with a "»" in front of it to clarify which choice is selected.
const TemplateDefault = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- if eq $.SelectedIndex $i }}
      {{- print "⇡\n"}}
      {{- print (Selected $choice) "\n" }}
    {{- else }}
      {{- print "⇡\n" }}
      {{- print (Unselected $choice) "\n" }}
    {{- end }}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- if eq $.SelectedIndex $i }}
      {{- print (Selected $choice) "\n" }}
      {{- print "⇣" }}
    {{- else }}
      {{- print (Unselected $choice) "\n" }}
      {{- print "⇣" }}
    {{- end }}
  {{- else -}}
    {{- if eq $.SelectedIndex $i }}
      {{- print (Selected $choice) "\n" }}
    {{- else }}
      {{- print (Unselected $choice) "\n" }}
    {{- end }}
  {{- end -}}
{{- end}}
`

// The template for displaying selection choices as table entries; shows the bolded prompt message followed on the next
// line by the filter prompt (if the prompt is using a filter) and then a blank line before the selection table. The
// next line is the header row for the table. The function to write the header row is given a boolean value: true if the
// table is paginated and the user can scroll up, otherwise false. After the header, the choices are rendered per their
// specified style. If the table is paginated and the user can scroll down, a "⇣" is inserted after the table.
const TemplateTable = `
{{- $canScrollUp := false -}}
{{- range $i, $choice := .Choices -}}
  {{- if IsScrollUpHintPosition $i -}}
    {{- $canScrollUp = true -}}
  {{- end -}}
{{- end -}}

{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{ range $i, $choice := .Choices }}
  {{- if (eq $i 0) }}{{ HeaderRow $canScrollUp }}{{ end }}
  {{- if IsScrollDownHintPosition $i -}}
    {{- if eq $.SelectedIndex $i }}
      {{- print (Selected $choice) "\n" }}
      {{- print "⇣" }}
    {{- else }}
      {{- print (Unselected $choice) "\n" }}
      {{- print "⇣" }}
    {{- end }}
  {{- else -}}
    {{- if eq $.SelectedIndex $i }}
      {{- print (Selected $choice) "\n" }}
    {{- else }}
      {{- print (Unselected $choice) "\n" }}
    {{- end }}
  {{- end -}}
{{- end}}
`

// The default template for displaying the result of the selection prompt, displaying the message followed by the choice
const ResultTemplateDefault = `
{{- print .Prompt " " (Final .FinalChoice) "\n" -}}
`

// A helper template for displaying the result of a selection prompt when the choices are structs. The name function
// is called to determine the string for displaying the choice instead of splatting the struct as a string.
const ResultTemplateByName = `
{{- print .Prompt " " (Final (Name .FinalChoice)) "\n" -}}
`

// This option overrides the default template the prompt uses to display to the terminal. For more information, see the
// docs for Selection: https://pkg.go.dev/github.com/erikgeiser/promptkit/selection#Selection.Template
func WithTemplate(template string) Option {
	return func(prompt *selection.Selection) {
		prompt.Template = template
	}
}

// This option overrides the default template the prompt uses to display results to the terminal. For more information,
// see the docs for Selection:
// https://pkg.go.dev/github.com/erikgeiser/promptkit/selection#Selection.ResultTemplate
func WithResultTemplate(template string) Option {
	return func(prompt *selection.Selection) {
		prompt.ResultTemplate = template
	}
}
