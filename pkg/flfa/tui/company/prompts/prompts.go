package prompts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	tprompts "github.com/FlagrantGarden/flfa/pkg/flfa/tui/traits/prompts"
	"github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	pterm "github.com/FlagrantGarden/flfa/pkg/tympan/printers/terminal"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/confirmer"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/selector"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/texter"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func ChooseCompany(settings *terminal.Settings, creating bool, companies []data.Company) *selection.Selection {
	var companyNames []string
	var message string

	if creating {
		companyNames = append(companyNames, "Create a new company")
		message = "Select an existing company to start from or create your own:"
	} else {
		message = "Select a company to edit:"
	}

	for _, company := range companies {
		companyNames = append(companyNames, company.Name)
	}

	return selector.NewStringSelector(
		message,
		companyNames,
		selector.WithPageSize(5),
		selector.WithSelectedChoiceStyle(selector.ColorizedBasicSelectedChoiceStyle(settings.ExtraColor("highlight"))),
	)
}

func ChooseCompanyModel(settings *terminal.Settings, creating bool, companies []data.Company) *selection.Model {
	return selection.NewModel(ChooseCompany(settings, creating, companies))
}

func GetName(settings *terminal.Settings) *textinput.TextInput {
	return texter.New(
		"What is this company called?",
		texter.WithPlaceholder("Name cannot be empty"),
	)
}

func GetNameModel(settings *terminal.Settings) *textinput.Model {
	return textinput.NewModel(GetName(settings))
}

func GetDescription(settings *terminal.Settings) *textinput.TextInput {
	return texter.New(
		"How would you describe this company?",
		texter.WithPlaceholder("Description cannot be empty"),
	)
}

func GetDescriptionModel(settings *terminal.Settings) *textinput.Model {
	return textinput.NewModel(GetDescription(settings))
}

func SelectOption(settings *terminal.Settings, canRemoveAGroup bool, hasCaptain bool) *selection.Selection {
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
		selector.WithSelectedChoiceStyle(selector.ColorizedBasicSelectedChoiceStyle(settings.ExtraColor("highlight"))),
	)
}

func SelectOptionModel(settings *terminal.Settings, canRemoveAGroup bool, hasCaptain bool) *selection.Model {
	return selection.NewModel(SelectOption(settings, canRemoveAGroup, hasCaptain))
}

type SelectGroupFor int

const (
	Editing SelectGroupFor = iota
	Copying
	Promoting
	Removing
)

func SelectGroup(settings *terminal.Settings, action SelectGroupFor, groups []data.Group) (prompt *selection.Selection) {
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
		tableSettings, _ := settings.Copy()
		tableSettings.SetFlagOn("for_selection")
		tableSettings.SetFlag("can_scroll_up", pterm.FlagFromBool(canScrollUp))

		return (&copy).TableHeaderTerminal(tableSettings)
	}

	selectedChoiceStyle := groupChoiceStyle(settings, true, action)

	unselectedChoiceStyle := groupChoiceStyle(settings, false, action)

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
		selector.WithPageSize(6),
	)
}

func SelectGroupModel(settings *terminal.Settings, action SelectGroupFor, groups []data.Group) *selection.Model {
	return selection.NewModel(SelectGroup(settings, action, groups))
}

func groupChoiceStyle(settings *terminal.Settings, selected bool, action SelectGroupFor) selector.ChoiceStyleFunc {
	return func(choice *selection.Choice) string {
		group, _ := choice.Value.(data.Group)
		promptSettings, _ := settings.Copy()

		if selected {
			promptSettings.SetFlagOn("selected")
		} else {
			promptSettings.SetFlagOff("selected")
		}

		if action == Removing {
			promptSettings.SetFlagOn("removing")
		}

		return group.ToTerminalTableEntry(promptSettings)
	}
}

func SelectCaptaincyOption(settings *terminal.Settings) *selection.Selection {
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
		selector.WithSelectedChoiceStyle(selector.ColorizedBasicSelectedChoiceStyle(settings.ExtraColor("highlight"))),
	)
}

func SelectCaptaincyOptionModel(settings *terminal.Settings) *selection.Model {
	return selection.NewModel(SelectCaptaincyOption(settings))
}

type CaptainRerollChoice struct {
	Message string
	Trait   data.Trait
}

func SelectRerollCaptainTrait(settings *terminal.Settings, group data.Group, availableTraits []data.Trait) *selection.Selection {
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
		messageStyle := settings.DynamicStyle("add_selected_trait_name").Width(maxWidth).PaddingRight(2)
		effectStyle := settings.DynamicStyle("selected_trait_effect").Width(120 - maxWidth)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			messageStyle.Render(fmt.Sprintf("» %s", info.Message)),
			effectStyle.Render(info.Trait.Effect),
		)
	}

	unselectedChoiceStyle := func(choice *selection.Choice) string {
		info, _ := choice.Value.(CaptainRerollChoice)
		messageStyle := settings.DynamicStyle("unselected_trait_name").Width(maxWidth).PaddingRight(2)
		effectStyle := settings.DynamicStyle("unselected_trait_effect").Width(120 - maxWidth)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			messageStyle.Render(fmt.Sprintf("  %s", info.Message)),
			effectStyle.Render(info.Trait.Effect),
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

func SelectRerollCaptainTraitModel(settings *terminal.Settings, group data.Group, availableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectRerollCaptainTrait(settings, group, availableTraits))
}

func SelectCaptainTrait(settings *terminal.Settings, availableTraits []data.Trait) *selection.Selection {
	return tprompts.SelectTrait(
		settings,
		"Which Trait do you want the Captain to have?",
		data.FilterTraitsByType("Captain", availableTraits),
	)
}

func SelectCaptainTraitModel(settings *terminal.Settings, availableTraits []data.Trait) *selection.Model {
	return selection.NewModel(SelectCaptainTrait(settings, availableTraits))
}

func ConfirmDemoteCaptain(settings *terminal.Settings, group data.Group) *confirmation.Confirmation {
	captain := settings.RenderWithDynamicStyle("warning_emphasis", group.Name)
	message := fmt.Sprintf("Are you sure you want to demote %s from being captain?", captain)

	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmDemoteCaptainModel(settings *terminal.Settings, group data.Group) *confirmation.Model {
	return confirmation.NewModel(ConfirmDemoteCaptain(settings, group))
}

func ConfirmReplaceCaptain(settings *terminal.Settings, current data.Group, new data.Group) *confirmation.Confirmation {
	currentCaptain := settings.RenderWithDynamicStyle("warning_emphasis", current.Name)
	newCaptain := settings.RenderWithDynamicStyle("confirmation_emphasis", new.Name)
	message := fmt.Sprintf(
		"Are you sure you want to demote %s from being captain and promote %s in their stead?",
		currentCaptain,
		newCaptain,
	)
	return confirmer.New(message, confirmer.WithDefaultValue(confirmation.No))
}

func ConfirmReplaceCaptainModel(settings *terminal.Settings, current data.Group, new data.Group) *confirmation.Model {
	return confirmation.NewModel(ConfirmReplaceCaptain(settings, current, new))
}
