package prompts

import (
	"regexp"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	pterm "github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/erikgeiser/promptkit/selection"
)

func SelectTrait(settings *terminal.Settings, message string, applicableTraits []data.Trait) *selection.Selection {
	filter := func(filter string, choice *selection.Choice) bool {
		chosenTrait, _ := choice.Value.(data.Trait)
		regex := regexp.MustCompile(strings.ToLower(filter))
		return regex.MatchString(strings.ToLower(chosenTrait.Name))
	}

	nameFunc := func(choice *selection.Choice) string {
		return choice.Value.(data.Trait).Name
	}

	longestTraitWidth := 0
	for _, trait := range applicableTraits {
		nameLength := len(trait.Name)
		if nameLength > longestTraitWidth {
			longestTraitWidth = nameLength
		}
	}
	longestTraitWidth += 1

	selectedChoiceStyle := TraitChoiceStyle(settings, true, longestTraitWidth)

	unselectedChoiceStyle := TraitChoiceStyle(settings, false, longestTraitWidth)

	return selector.NewStructSelector(
		message,
		applicableTraits,
		filter,
		nameFunc,
		selector.WithSelectedChoiceStyle(selectedChoiceStyle),
		selector.WithUnselectedChoiceStyle(unselectedChoiceStyle),
		selector.WithPageSize(5),
	)
}

func SelectTraitModel(settings *terminal.Settings, message string, applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectTrait(settings, message, applicableTraits))
}

func TraitChoiceStyle(settings *terminal.Settings, selected bool, leadWidth int, options ...pterm.Option) func(choice *selection.Choice) string {
	return func(choice *selection.Choice) string {
		trait, _ := choice.Value.(data.Trait)
		choiceSettings, _ := settings.Copy()
		choiceSettings.SetFlag("selected", pterm.FlagFromBool(selected))
		return trait.ToTerminalChoice(choiceSettings, leadWidth)
	}
}
