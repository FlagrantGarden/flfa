package dynamic

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

// The Model for a dynamic prompt is able to display confirmation, selection, and text input prompts from data without
// the developer knowing ahead of time which type of prompt will be needed. This allows a program to take advantage of
// prompts defined in data outside of the go code itself.
type Model struct {
	// Which prompt is active
	ActiveType Prompt
	// The model for the selection prompt
	Selection *selection.Model
	// The model for the text input prompt
	TextInput *textinput.Model
	// The model for the confirmation prompt
	Confirmation *confirmation.Model
}

// When the Model is initialized, it passes straight through to appropriate submodel for the prompt.
func (model *Model) Init() tea.Cmd {
	switch model.ActiveType {
	case Confirmation:
		return model.Confirmation.Init()
	case Selection:
		return model.Selection.Init()
	case TextInput:
		return model.TextInput.Init()
	default:
		return nil
	}
}

// When the Model is viewed, it passes straight through to the appropriate submodel for the prompt. If no prompt is
// active or the prompt type could not be determined, it displays an error.
func (model *Model) View() string {
	switch model.ActiveType {
	case Confirmation:
		return model.Confirmation.View()
	case Selection:
		return model.Selection.View()
	case TextInput:
		return model.TextInput.View()
	default:
		return fmt.Sprintf("Error: dynamic prompt (type %s) is not a valid type for active display", model.ActiveType)
	}
}

// When the Model is updated, it passes straight through to the appropriate submodel for the prompt.
func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch model.ActiveType {
	case Confirmation:
		return model.Confirmation.Update(msg)
	case Selection:
		return model.Selection.Update(msg)
	case TextInput:
		return model.TextInput.Update(msg)
	default:
		return model, nil
	}
}

// The New function takes the data for dynamic prompt info and returns a dynamic prompt model configured per the values
// and options defined in the Info object.
func New(info Info) *Model {
	model := &Model{}

	switch info.EnumType() {
	case Confirmation:
		model.ActiveType = Confirmation

		var options []confirmer.Option
		for _, option := range info.Options {
			options = append(options, option.ConfirmationOptions()...)
		}

		prompt := confirmer.New(info.Message, options...)
		model.Confirmation = confirmation.NewModel(prompt)
	case Selection:
		model.ActiveType = Selection

		var options []selector.Option
		for _, option := range info.Options {
			options = append(options, option.SelectionOptions()...)
		}

		prompt := selector.New(info.Message, selection.Choices([]string{}), options...)
		model.Selection = selection.NewModel(prompt)
	case TextInput:
		model.ActiveType = TextInput

		var options []texter.Option
		for _, option := range info.Options {
			options = append(options, option.TextInputOptions()...)
		}

		prompt := texter.New(info.Message, options...)
		model.TextInput = textinput.NewModel(prompt)
	}

	return model
}
