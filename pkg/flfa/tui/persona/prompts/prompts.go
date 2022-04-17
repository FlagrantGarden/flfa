package prompts

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/state/player"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func GetName() *textinput.TextInput {
	return texter.New(
		"What name do you want to be called by?",
		texter.WithPlaceholder("Name cannot be empty."),
		texter.WithInputWidth(30),
	)
}

func GetNameModel() *textinput.Model {
	return textinput.NewModel(GetName())
}

func Choose(personas []player.Player) *selection.Selection {
	var messageBuilder strings.Builder
	messageBuilder.WriteString("It looks like this is your first time playing ")
	messageBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true).Render("Flagrant Factions"))
	messageBuilder.WriteString("!\nDo you want to create a user persona now? ")
	messageBuilder.WriteString("If not, you can just use the default one.\n")
	message := messageBuilder.String()

	var choices []string
	choices = append(choices, "Create a new persona")
	if len(personas) > 0 {
		message = "Which persona would you like to use?"
		for _, persona := range personas {
			choices = append(choices, persona.Name)
		}
	}

	if !utils.Contains(choices, "default") {
		choices = append(choices, "default")
	}

	return selector.NewStringSelector(message, choices, selector.WithPageSize(5))
}

func ChoosePersonaModel(personas []player.Player) *selection.Model {
	return selection.NewModel(Choose(personas))
}

func SetAsPreferred(name string) *confirmation.Confirmation {
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("Do you want to set %s as your preferred persona? ", name))
	messageBuilder.WriteString("Next time you play, this one will load automatically.")

	return confirmer.New(messageBuilder.String(), confirmer.WithDefaultValue(confirmation.No))
}

func SetAsPreferredModel(name string) *confirmation.Model {
	return confirmation.NewModel(SetAsPreferred(name))
}
