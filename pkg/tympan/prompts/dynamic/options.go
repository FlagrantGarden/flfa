package dynamic

import (
	"strings"
	"text/template"

	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
)

// The Option enum maps an option string to a defined behavior for the dynamic prompt.
type Option int

const (
	// The specified string did not map to a valid option
	OptionUnknown Option = iota
	// The option sets the default value for a confirmation prompt to yes, no, or undecided.
	ConfirmationDefault
	// The option adds the value list to the selection prompt as valid choices. Use this when the choices are a list of
	// non-complex objects, like strings or integers.
	SelectionChoiceSimple
	// The option adds the value list to the selection prompt as valid choices. Use this when the choices are a list of
	// complex objects as maps; to use this option, each choice in the data must include a name key with a string value
	// for identifying the choice.
	SelectionChoiceComplex
)

// Returns the string representation of an Option enum
func (option Option) String() (value string) {
	switch option {
	case OptionUnknown:
		value = "OptionUnkown"
	case ConfirmationDefault:
		value = "ConfirmationDefault"
	case SelectionChoiceSimple:
		value = "SelectionChoiceSimple"
	case SelectionChoiceComplex:
		value = "SelectionChoiceComplex"
	}
	return value
}

// Prompt Options must declare a type and a value. They are used for changing the behavior of a dynamic prompt.
type PromptOption struct {
	// The Type must be a string that maps to a valid Option via the EnumType() method.
	Type string
	// The Value can be anything but is most commonly a string, a map[string]any, or a slice of either.
	Value any
}

// The EnumType method maps the string specified in the data for a dynamic prompt's options to an Option enum. If the
// specified string cannot be mapped, it returns OptionUnknown and is ignored.
func (option PromptOption) EnumType() (modifier Option) {
	switch strings.ToLower(option.Type) {
	case "confirmation_default":
		modifier = ConfirmationDefault
	case "selection_choices_simple":
		modifier = SelectionChoiceSimple
	case "selection_choices_complex":
		modifier = SelectionChoiceComplex
	}

	return modifier
}

// The SelectionOptions method introspects on a PromptOption to return a slice of options to use when creating a
// selection prompt.
//
// If the prompt option is for a simple selection choice, it appends the value of the prompt option to the list of
// valid choices.
//
// If the prompt option is for a complex selection choice, it returns options which:
//
// 1. Append the slice of choices under the "choices" key in the option's value as valid choices
// 2. Add a name function to the prompt's extended template funcs, returning the value of the name key (as specified
// in the "name" key of the option's value) from each choice when rendering the choice in the prompt
// 3. Add a filter function to enable users to type to filter the available choices, using the name key again
// 4. Specifies the alternate result template to ensure the result is displayed using the name and not a splat of the
// choice's values as a string.
func (option PromptOption) SelectionOptions() (options []selector.Option) {
	switch option.EnumType() {
	case SelectionChoiceSimple:
		options = append(options, selector.WithChoices(option.Value))
	case SelectionChoiceComplex:
		value := option.Value.(map[string]any)
		nameKey := value["name"].(string)

		options = append(options, selector.WithChoices(value["choices"]))

		options = append(options, selector.WithExtendedTemplateFuncs(template.FuncMap{
			"name": func(choice *selection.Choice) string {
				return choice.Value.(map[string]any)[nameKey].(string)
			},
		}))

		options = append(options, selector.WithFilter(func(filter string, choice *selection.Choice) bool {
			name := choice.Value.(map[string]any)[nameKey].(string)
			return strings.Contains(strings.ToLower(name), strings.ToLower(filter))
		}))

		options = append(options, selector.WithResultTemplate(selector.ResultTemplateByName))
	}

	return options
}

// The TextInputOptions method introspects on a PromptOption to return a slice of options to use when creating a
// textinput prompt.
//
// No options are yet implemented.
func (option PromptOption) TextInputOptions() (options []texter.Option) {
	// none implemented yet
	return options
}

// The ConfirmationOptions method introspects on a PromptOption to return a slice of options to use when creating a
// confirmation prompt.
//
// If the prompt option is for a confirmation default, it sets the prompt's default value to Yes (if the value in the
// data is true or "yes"), No (if the value in the data is false or "no"), or Undecided (if the value is anything else).
func (option PromptOption) ConfirmationOptions() (options []confirmer.Option) {
	switch option.EnumType() {
	case ConfirmationDefault:
		var defaultValue confirmation.Value

		switch option.Value {
		case true, "yes":
			defaultValue = confirmation.Yes
		case false, "no":
			defaultValue = confirmation.No
		default:
			defaultValue = confirmation.Undecided
		}

		options = append(options, confirmer.WithDefaultValue(defaultValue))
	}
	return options
}
