package dynamic

import (
	"strings"
)

// The Prompt enum declares the types of prompts that the dynamic model can handle.
type Prompt int

const (
	// Unknown prompts go unhandled, erroring.
	PromptUnknown Prompt = iota
	// A confirmation prompt asks a user a yes/no question
	Confirmation
	// A selection prompt has the user select one choice from a list of choices
	Selection
	// A text input prompt has the user enter a string of text
	TextInput
)

// Returns the string representation of a Prompt enum
func (prompt Prompt) String() (value string) {
	switch prompt {
	case PromptUnknown:
		value = "PromptUnkown"
	case Confirmation:
		value = "Confirmation"
	case Selection:
		value = "Selection"
	case TextInput:
		value = "TextInput"
	}
	return value
}

// An Info object is used to dynamically build a prompt from data, enabling the creation of prompts defined outside of
// the compiled code itself.
type Info struct {
	// The Type represents a valid prompt and must be mappable to the Prompt enum via the EnumType() method.
	Type string
	// The Message is an arbitrary string that will be used for the prompt message.
	Message string
	// Dynamic prompts can include zero or more options which alter the behavior of the prompt.
	Options []PromptOption
}

// The EnumType() method returns the valid prompt type that the info maps to. If the specified string does not map to
// any valid prompt type, it returns PromptUnknown. This method is case insensitive.
func (info Info) EnumType() (prompt Prompt) {
	switch strings.ToLower(info.Type) {
	case "confirmation":
		prompt = Confirmation
	case "selection":
		prompt = Selection
	case "text", "textinfo":
		prompt = TextInput
	}

	return prompt
}
