package prompts

import (
	"regexp"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	pterm "github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/erikgeiser/promptkit/selection"
)

func SelectTrait(message string, applicableTraits []data.Trait) *selection.Selection {
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

	selectedChoiceStyle := TraitChoiceStyle(true, longestTraitWidth)

	unselectedChoiceStyle := TraitChoiceStyle(false, longestTraitWidth)

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

func SelectTraitModel(message string, applicableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectTrait(message, applicableTraits))
}

func TraitChoiceStyle(selected bool, leadWidth int, options ...pterm.Option) func(choice *selection.Choice) string {
	return func(choice *selection.Choice) string {
		trait, _ := choice.Value.(data.Trait)
		return trait.ToTerminalChoice(selected, leadWidth, options...)
	}
}
