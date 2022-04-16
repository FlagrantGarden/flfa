package prompts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	tprompts "github.com/FlagrantGarden/flfa/pkg/flfa/tui/traits/prompts"
	pterm "github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func ChooseCompany(companies []data.Company) *selection.Selection {
	var companyNames []string
	companyNames = append(companyNames, "Create a new company")
	for _, company := range companies {
		companyNames = append(companyNames, company.Name)
	}

	return selector.NewStringSelector(
		"Which company would you like to play as?",
		companyNames,
		selector.WithPageSize(5),
	)
}

func ChooseCompanyModel(companies []data.Company) *selection.Model {
	return selection.NewModel(ChooseCompany(companies))
}

func GetName() *textinput.TextInput {
	return texter.New(
		"What is your company called?",
		texter.WithPlaceholder("Name cannot be empty"),
	)
}

func GetNameModel() *textinput.Model {
	return textinput.NewModel(GetName())
}

func GetDescription() *textinput.TextInput {
	return texter.New(
		"How would you describe your company?",
		texter.WithPlaceholder("Description cannot be empty"),
	)
}

func GetDescriptionModel() *textinput.Model {
	return textinput.NewModel(GetDescription())
}

func SelectOption(canRemoveAGroup bool, hasCaptain bool) *selection.Selection {
	options := []string{
		"Save & Continue",
		"Change Name",
		"Change Description",
		"Create & add a new Group",
		"Add a copy of a Group",
		"Edit a Group",
	}

	if canRemoveAGroup {
		options = append(options, "Remove a Group")
	}

	if hasCaptain {
		options = append(options, "Update Captaincy")
	} else {
		options = append(options, "Promote a Group to Captain")
	}

	return selector.NewStringSelector(
		"What would you like to do with this Company?",
		options,
		selector.WithPageSize(5),
	)
}

func SelectOptionModel(canRemoveAGroup bool, hasCaptain bool) *selection.Model {
	return selection.NewModel(SelectOption(canRemoveAGroup, hasCaptain))
}

type SelectGroupFor int

const (
	Editing SelectGroupFor = iota
	Copying
	Promoting
	Removing
)

func SelectGroup(action SelectGroupFor, groups []data.Group) (prompt *selection.Selection) {
	filter := func(filter string, choice *selection.Choice) bool {
		chosenGroup, _ := choice.Value.(data.Group)
		regex := regexp.MustCompile(strings.ToLower(filter))
		return regex.MatchString(strings.ToLower(chosenGroup.Name))
	}

	nameFunc := func(choice *selection.Choice) string {
		return choice.Value.(data.Group).Name
	}

	headerRowFunc := func(canScrollUp bool) string {
		copy := groups[0]

		return (&copy).TableHeaderTerminal(
			pterm.WithFlagOn("for_selection"),
			pterm.WithFlag("can_scroll_up", pterm.FlagFromBool(canScrollUp)),
		)
	}

	selectedChoiceStyle := groupChoiceStyle(true, action)

	unselectedChoiceStyle := groupChoiceStyle(false, action)

	var message strings.Builder
	switch action {
	case Copying:
		message.WriteString("Which Group would you like to make a copy of?\n")
	case Editing:
		message.WriteString("Which Group would you like to edit?\n")
	case Promoting:
		message.WriteString("Which Group would you like to promote to captain?\n")
	case Removing:
		message.WriteString("Which Group would you like to remove?\n")
	}

	return selector.NewTableSelector(
		message.String(),
		groups,
		filter,
		nameFunc,
		headerRowFunc,
		selectedChoiceStyle,
		unselectedChoiceStyle,
		selector.WithPageSize(3),
	)
}

func SelectGroupModel(action SelectGroupFor, groups []data.Group) *selection.Model {
	return selection.NewModel(SelectGroup(action, groups))
}

func groupChoiceStyle(selected bool, action SelectGroupFor) selector.ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		group, _ := choice.Value.(data.Group)
		var options []pterm.Option

		if selected {
			options = append(options, pterm.WithFlagOn("selected"))
		} else {
			options = append(options, pterm.WithFlagOff("selected"))
		}

		var style lipgloss.Style
		switch action {
		case Editing:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("32"))
		case Removing:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("166"))
		}
		options = append(options, pterm.WithExtraStyle("selected_lead", style))

		return group.ToTerminalTableEntry(options...)
	}
}

func SelectCaptaincyOption() *selection.Selection {
	options := []string{
		"Go back",
		"Reroll Captain's trait",
		"Choose Captain's trait",
		"Demote Captain",
		"Choose a different Captain",
	}

	return selector.NewStringSelector(
		"What do you want to do about the Captain?",
		options,
		selector.WithPageSize(5),
	)
}

func SelectCaptaincyOptionModel() *selection.Model {
	return selection.NewModel(SelectCaptaincyOption())
}

type CaptainRerollChoice struct {
	Message string
	Trait   data.Trait
}

func SelectRerollCaptainTrait(group data.Group, availableTraits []data.Trait) *selection.Selection {
	copy := group
	copy.PromoteToCaptain(nil, availableTraits...)
	options := []CaptainRerollChoice{
		{
			Message: fmt.Sprintf("Keep the new trait (%s) & return", copy.Captain.Name),
			Trait:   copy.Captain,
		},
		{
			Message: fmt.Sprintf("Keep the old trait (%s) & return", group.Captain.Name),
			Trait:   group.Captain,
		},
		{
			Message: "Roll again",
		},
	}

	maxWidth := 35
	for _, option := range options {
		length := len(option.Message)
		if length > maxWidth {
			maxWidth = length
		}
	}
	maxWidth += 5

	filter := func(filter string, choice *selection.Choice) bool {
		chosen, _ := choice.Value.(CaptainRerollChoice)
		regex := regexp.MustCompile(strings.ToLower(filter))
		return regex.MatchString(strings.ToLower(chosen.Message))
	}

	nameFunc := func(choice *selection.Choice) string {
		return choice.Value.(CaptainRerollChoice).Message
	}

	selectedChoiceStyle := func(choice *selection.Choice) string {
		info, _ := choice.Value.(CaptainRerollChoice)
		messageStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32")).Width(maxWidth).PaddingRight(2)
		effectStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Width(120 - maxWidth)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			messageStyle.Render(fmt.Sprintf("» %s", info.Message)),
			effectStyle.Render(info.Trait.Effect),
		)
	}

	unselectedChoiceStyle := func(choice *selection.Choice) string {
		info, _ := choice.Value.(CaptainRerollChoice)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(maxWidth).PaddingRight(2).Render(fmt.Sprintf("  %s", info.Message)),
			lipgloss.NewStyle().Width(120-2-maxWidth).Faint(true).Render(info.Trait.Effect),
		)
	}

	return selector.NewStructSelector(
		"What do you want to do?",
		options,
		filter,
		nameFunc,
		selector.WithPageSize(5),
		selector.WithSelectedChoiceStyle(selectedChoiceStyle),
		selector.WithUnselectedChoiceStyle(unselectedChoiceStyle),
	)
}

func SelectRerollCaptainTraitModel(group data.Group, availableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectRerollCaptainTrait(group, availableTraits))
}

func SelectCaptainTrait(availableTraits []data.Trait) *selection.Selection {
	return tprompts.SelectTrait(
		"Which Trait do you want the Captain to have?",
		data.FilterTraitsByType("Captain", availableTraits),
	)
}

func SelectCaptainTraitModel(availableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectCaptainTrait(availableTraits))
}

func ConfirmDemoteCaptain(group data.Group) *confirmation.Confirmation {
	captain := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("166")).Render(group.Name)
	message := fmt.Sprintf("Are you sure you want to demote %s from being captain?", captain)

	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmDemoteCaptainModel(group data.Group) *confirmation.Model {
	return confirmation.NewModel(ConfirmDemoteCaptain(group))
}

func ConfirmReplaceCaptain(current data.Group, new data.Group) *confirmation.Confirmation {
	style := lipgloss.NewStyle().Bold(true)
	currentCaptain := style.Copy().Foreground(lipgloss.Color("166")).Render(current.Name)
	newCaptain := style.Copy().Foreground(lipgloss.Color("32")).Render(new.Name)
	message := fmt.Sprintf(
		"Are you sure you want to demote %s from being captain and promote %s in their stead?",
		currentCaptain,
		newCaptain,
	)
	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmReplaceCaptainModel(current data.Group, new data.Group) *confirmation.Model {
	return confirmation.NewModel(ConfirmReplaceCaptain(current, new))
}
