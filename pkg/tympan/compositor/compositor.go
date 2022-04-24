package compositor

import (
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

// The State of a Compositor is used extensively when determining how to process view and update calls to the model.
// The shared states defined in this library can be used by any model inheriting the Compositor. When defining
// additional states specific to your application, be sure to start your list of State constants with a higher integer.
//
// For example:
//
//     const (
//       StateCreating compositor.State = iota + 100
//       StateEditing
//     )
//
// The States you define should be unique within your application and not conflict with the states defined here.
type State int

const (
	// When the state is unknown, something has gone wrong but the model is not necessarily broken
	StateUnknown State = iota
	// A fatal error has occurred
	StateBroken
	// The model is not actively in use but has been initialized and may be used again
	StateReady
	// The model has ended, cancelled by the user
	StateCancelled
	// The model has ended, the user has chosen to return to the program flow
	StateDone
	// The model is paused while the application saves
	StateSavingConfiguration
	// The model has finished saving the application but is not yet in another state
	StateSavedConfiguration
)

// A Compositor is a shareable configuration for TUIs built with bubbletea, promptkit, and lipgloss. It provides some
// default behaviors and expectations for how a TUI is constructed and used.
type Compositor struct {
	// The state of a compositor-based TUI is used to determine how to flow updates and views.
	State State
	// Compositors may be used in models at any level of a TUI; when compositors are submodels, they need to be able to
	// handle reporting that they are no longer active and why.
	IsSubmodel bool
	// Compositors can record and display fatal errors to the user; this simplifies handling non-recoverable errors.
	FatalError error
	// Compositors can be instantiated with a particular width and update it as needed. This field caches that value so
	// it can be reused with querying the terminal.
	Width int
	// Compositors can be instantiated with a particular height and update it as needed. This field caches that value so
	// it can be reused with querying the terminal.
	Height int
	// Compositors always include a confirmation model, but it need not be used.
	Confirmation *confirmation.Model
	// Compositors always include a confirmation model, but it need not be used.
	Selection *selection.Model
	// Compositors always include a confirmation model, but it need not be used.
	TextInput *textinput.Model
	// Compositors always include default terminal settings; they can be extended or replaced.
	TerminalSettings *terminal.Settings
}

// A Modeler is a bubbletea Model which follows the design pattern the Compositor supports; namely the use of state and
// substate to control flow and view of the TUI.
type Modeler interface {
	// The Init method runs when the TUI starts
	Init() tea.Cmd
	// The SetAndStartState method is used to move into another top-level state and typically that state's default
	// substate.
	SetAndStartState(state State) (cmd tea.Cmd)
	// The Update method runs whenever the TUI receives input from the user or system and should explicitly pass the
	// update message to UpdateOnKeyPress (if the user pressed a key), UpdateOnSubmodelEnded (if the TUI received an end
	// message from a submodel), or UpdateFallThrough (if the current model does not explicitly handle the input message).
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	// The UpdateFallThrough method passes the update message to the correct submodel based on the current state/substate.
	UpdateFallThrough(msg tea.Msg) (cmd tea.Cmd)
	// The UpdateOnKeyPress method handles the update message based on the current state/substate and key pressed.
	UpdateOnKeyPress(msg tea.KeyMsg) (cmd tea.Cmd)
	// The UpdateOnSubmodelEnded method handles what to do when a submodel ends based on the current state/substate.
	UpdateOnSubmodelEnded() (cmd tea.Cmd)
	// The View method returns a string for rendering to the terminal
	View() string
}

// Compositor-based Models should use the variadic option model by having their New function accept options to change
// the model before returning it. These options should be defined as functions with semantically clear names and only
// modify the model. For example:
//
//     type MyModel struct {
//       *compositor.Compositor
//       Required string
//       Foo string
//     }
//     func WithFoo(foo string) compositor.Option[*MyModel] {
//       return func(model *MyModel) {
//         model.Foo = foo
//       }
//     }
//     func New(required, options ...compositor.Option[*MyModel]) *MyModel {
//       model := &MyModel{Required: required}
//
//       for _, option := range options {
//         option(model)
//       }
//
//       return model
//     }
type Option[Model Modeler] func(model Model)

// Compositor-based Models should use substates if they have multiple distinct flows depending on the higher-level
// state; for example, if the TUI creates an object and edits it, it might have the creating and editing substates.
// Substates provide a standardized interface for managing flow and view of the model.
type SubstateInterface[Model Modeler] interface {
	// The Start method for a substate should perform the necessary setup for the TUI to move into this substate
	Start(model Model) (cmd tea.Cmd)
	// The UpdateOnEnded method for a substate handles when a submodel associated with the substate reports an
	// end message
	UpdateOnEnded(model Model) (cmd tea.Cmd)
	// The UpdateOnEnter method for a substate handles when the user presses the enter key in that substate
	UpdateOnEnter(model Model) (cmd tea.Cmd)
	// The UpdateOnEsc method for a substate handles when the user presses the escape key in that substate
	UpdateOnEsc(model Model) (cmd tea.Cmd)
	// The UpdateOnFallThrough method passes the update message to the correct submodel
	UpdateOnFallThrough(model Model, msg tea.Msg) (cmd tea.Cmd)
	// The View method ensures the correct view for this substate is rendered and returned
	View(model Model) (view string)
}

// When a compositor-based model is used as a submodel for another, it should report an end message instead of quitting
// the TUI.
type EndMsg struct{}

// When a compositor-based model encounters an error it cannot handle, it should record that error as fatal; this method
// is a shorthand for doing so, saving the error and reporting that the model is broken.
func (model *Compositor) RecordFatalError(err error) tea.Cmd {
	model.FatalError = err
	return model.Broken
}

// This command should be called when an unhandleable problem occurs in the model
func (model *Compositor) Broken() tea.Msg {
	model.State = StateBroken
	return EndMsg{}
}

// This command should be called when the user chooses to exit early from a compositor-based model when it is being used
// as a submodel.
func (model *Compositor) Cancelled() tea.Msg {
	model.State = StateCancelled
	if model.IsSubmodel {
		return EndMsg{}
	} else {
		return tea.Quit()
	}
}

// This command should be called when the user chooses to exit gracefully after completing a compositor-based model when
// it is being used as a submodel.
func (model *Compositor) Done() tea.Msg {
	model.State = StateDone
	if model.IsSubmodel {
		return EndMsg{}
	} else {
		return tea.Quit()
	}
}

// This method sets the cached width and height of a compositor-based model but does not otherwise change any behavior.
func (sharedModel *Compositor) SetSize(width int, height int) {
	sharedModel.Width = width
	sharedModel.Height = height
}

// This method returns the default view for a compositor-based model when a fatal error has been recorded.
func (sharedModel *Compositor) ViewFatalError(options ...terminal.Option) string {
	message := lipgloss.JoinVertical(
		lipgloss.Center,
		sharedModel.TerminalSettings.DynamicStyle("error_header").Render("Fatal Error!"),
		sharedModel.TerminalSettings.DynamicStyle("error_message").Render(sharedModel.FatalError.Error()),
	)

	message = sharedModel.TerminalSettings.DynamicStyle("error_box").Render(message)
	// message = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(errColor).Render(message)
	return lipgloss.Place(120, 60, lipgloss.Center, lipgloss.Center, message)
}

// Returns the default terminal settings required for the Compositor to behave and display correctly.
func DefaultTerminalSettings() *terminal.Settings {
	return terminal.New(
		terminal.WithExtraColor("error", lipgloss.Color("166")),
		terminal.WithExtraStyle("error_header", lipgloss.NewStyle().Bold(true).Margin(1)),
		terminal.WithExtraStyle("error_message", lipgloss.NewStyle().Align(lipgloss.Left).Width(80).Padding(0, 2, 1)),
		terminal.WithExtraStyle("error_box", lipgloss.NewStyle().Border(lipgloss.DoubleBorder())),
		terminal.WithDynamicStyle(
			"error_header",
			terminal.OverrideWithExtraStyle("error_header"),
			terminal.ColorizeForeground("error"),
		),
		terminal.WithDynamicStyle(
			"error_message",
			terminal.OverrideWithExtraStyle("error_message"),
			terminal.ColorizeForeground("error"),
		),
		terminal.WithDynamicStyle(
			"error_box",
			terminal.OverrideWithExtraStyle("error_box"),
			terminal.ColorizeBorderForeground("error"),
		),
	)
}

// Returns a new instance of a Compositor with the terminal settings initialized. Specify one or more options to append
// or override the existing settings, if desired.
func New(terminalSettingsOptions ...terminal.Option) *Compositor {
	var combinedOptions []terminal.Option
	combinedOptions = append(combinedOptions, terminal.From(*DefaultTerminalSettings()))
	combinedOptions = append(combinedOptions, terminalSettingsOptions...)

	return &Compositor{
		TerminalSettings: terminal.New(combinedOptions...),
	}
}
