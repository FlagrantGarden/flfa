package data

import (
	// "github.com/FlagrantGarden/flfa/pkg/flfa"

	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/FlagrantGarden/flfa/pkg/tympan/module/scripting"
	"github.com/FlagrantGarden/flfa/pkg/tympan/prompts/dynamic"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/justinian/dice"
)

type Trait struct {
	Name      string
	Type      string
	Source    string
	Roll      int
	Effect    string
	Points    int
	Scripting TraitScripting
	Choices   []*TraitChoice
}

type TraitScripting struct {
	Requirements []string
	OnAdd        []string               `mapstructure:"on_add"`
	OnRemove     []string               `mapstructure:"on_remove"`
	InPlay       []TraitScriptingInPlay `mapstructure:"in_play"`
}

type TraitChoice struct {
	Name   string
	Value  any
	Prompt dynamic.Info // TODO: support serial prompt, nested prompts
}

type TraitScriptingInPlay struct {
	RegisterFor []string `mapstructure:"register_for"`
	AppliesTo   string   `mapstructure:"applies_to"`
	Uses        []TraitUses
	When        []string
	Then        []string
}

type TraitUses struct {
	PerTurn       int `mapstructure:"per_turn"`
	GlobalPerTurn int `mapstructure:"global_per_turn"`
}

func (trait Trait) WithSource(source string) Trait {
	trait.Source = source
	return trait
}

func (trait Trait) WithSubtype(subtype string) Trait {
	trait.Type = subtype
	return trait
}

type ApplicableRequirement func(group Group, trait Trait) bool

func WithGroupMaxPoints(max int) ApplicableRequirement {
	return func(group Group, trait Trait) bool {
		return group.Points+trait.Points <= max
	}
}

func WithCompanyMaxPoints(max int, current int) ApplicableRequirement {
	return func(group Group, trait Trait) bool {
		return current+trait.Points <= max
	}
}

func (trait Trait) Applicable(group Group, baseProfile Group, engine *scripting.Engine, requirements ...ApplicableRequirement) (bool, error) {
	// If the group already has the trait, it's definitely not applicable
	if utils.Contains(group.Traits, trait.Name) {
		return false, nil
	}

	// If the modified point cost would bring their points above 12 or below 1, unapplicable
	if group.Points+trait.Points < 1 {
		return false, nil
	}

	// Check additional requirements
	for _, requirement := range requirements {
		if !requirement(group, trait) {
			return false, nil
		}
	}

	errorPrefix := "unable to check if trait '%s' is applicable to the '%s' group:"
	name := fmt.Sprintf("CheckIfTraitApplicable: '%s'", trait.Name)
	script := engine.GetScript(name)
	if script == nil {
		body := trait.RequirementsScriptBody()
		err := engine.AddScript(name, body)
		if err != nil {
			return false, fmt.Errorf("%s %s", errorPrefix, err)
		}
		script = engine.GetScript(name)
	}

	tengoizedGroup, err := scripting.ConvertToTengoMap(group)
	if err != nil {
		return false, fmt.Errorf("%s %s", errorPrefix, err)
	}
	script.Add("profile", tengoizedGroup)

	// it's possible the base profile should just be stored on the model in tengoized form
	// it should never be modified, only replaced if the base profile is updated.
	tengoizedBaseProfile, err := scripting.ConvertToTengoMap(baseProfile)
	if err != nil {
		return false, fmt.Errorf("%s %s", errorPrefix, err)
	}
	script.Add("base_profile", tengoizedBaseProfile)

	result, err := script.Run()
	if err != nil {
		return false, fmt.Errorf("%s %s", errorPrefix, err)
	}

	return result.Get("trait_requirements_met").Bool(), nil
}

func (trait Trait) RequirementsScriptBody() string {
	var scriptBuilder strings.Builder
	// scriptBuilder.WriteString("profile := \"unset\"\n")
	// scriptBuilder.WriteString("base_profile := \"unset\"\n")
	scriptBuilder.WriteString("trait_requirements_met := true\n")
	for _, requirement := range trait.Scripting.Requirements {
		scriptBuilder.WriteString("if trait_requirements_met == true {\n")
		scriptBuilder.WriteString(fmt.Sprintf("  trait_requirements_met = %s\n", requirement))
		scriptBuilder.WriteString("}\n")
	}
	return scriptBuilder.String()
}

func (trait Trait) TraitWithChoiceUpdatedName() *Trait {
	for _, choice := range trait.Choices {
		regex := regexp.MustCompile(fmt.Sprintf("\\[%s\\]", choice.Name))
		trait.Name = regex.ReplaceAllLiteralString(trait.Name, choice.Value.(string))
	}
	return &trait
}

func (trait Trait) AddToGroup(group *Group, engine *scripting.Engine) (updatedGroup *Group, err error) {
	errorPrefix := fmt.Sprintf("can't add trait '%s' to Group '%s'", trait.Name, group.Name)
	// Some choices change the trait name, like [Kind]bane -> Bearbane
	for _, choice := range trait.Choices {
		regex := regexp.MustCompile(fmt.Sprintf("\\[%s\\]", choice.Name))
		trait.Name = regex.ReplaceAllLiteralString(trait.Name, choice.Value.(string))
	}
	// If the group already has the trait, bail out
	if utils.Contains(group.Traits, trait.Name) {
		return group, fmt.Errorf("%s: the group already has it.", errorPrefix)
	}

	name := fmt.Sprintf("AddTraitToGroup: '%s'", trait.Name)
	script := engine.GetScript(name)
	if script == nil {
		body := trait.OnAddScriptBody()
		// If there was nothing to do, don't run any scripts
		if body == "" {
			group.Traits = append(group.Traits, trait.Name)
			group.Points += trait.Points
			return group, nil
		}
		err := engine.AddScript(name, body)
		if err != nil {
			return group, fmt.Errorf("%s: %s", errorPrefix, err)
		}
		script = engine.GetScript(name)
	}

	tengoizedGroup, err := scripting.ConvertToTengoMap(group)
	if err != nil {
		return group, fmt.Errorf("%s: %s", errorPrefix, err)
	}
	script.Add("group", tengoizedGroup)

	tengoizedChoices := make(map[string]any)
	for _, choice := range trait.Choices {
		tengoizedChoices[choice.Name] = choice.Value
	}
	script.Add("choices", tengoizedChoices)

	result, err := script.Run()
	if err != nil {
		groupString := fmt.Sprintf("Tengoized Group: %s", tengoizedGroup)
		return group, fmt.Errorf("%s: %s\n%s", errorPrefix, err, groupString)
	}
	output_group, err := scripting.ConvertFromTengoMap[Group](result.Get("group").Map())
	if err != nil {
		return group, fmt.Errorf("%s: %s", errorPrefix, err)
	}
	updatedGroup = &output_group
	updatedGroup.Traits = append(updatedGroup.Traits, trait.Name)
	updatedGroup.Points += trait.Points
	return updatedGroup, nil
}

func (trait Trait) OnAddScriptBody() string {
	var scriptBuilder strings.Builder
	for _, change := range trait.Scripting.OnAdd {
		scriptBuilder.WriteString(fmt.Sprintf("%s\n", change))
	}
	return scriptBuilder.String()
}

func (trait Trait) RemoveFromGroup(group *Group, engine *scripting.Engine) (updatedGroup *Group, err error) {
	errorPrefix := fmt.Sprintf("can't remove trait '%s' to Group '%s'", trait.Name, group.Name)
	// If the group doesn't have the trait, bail out
	if !utils.Contains(group.Traits, trait.Name) {
		return group, fmt.Errorf("%s: the group doesn't have it.", errorPrefix)
	}

	name := fmt.Sprintf("RemoveTraitFromGroup: '%s'", trait.Name)
	script := engine.GetScript(name)
	if script == nil {
		body := trait.OnRemoveScriptBody()
		// If there was nothing to do, don't run any scripts
		if body == "" {
			for index, groupTrait := range group.Traits {
				if groupTrait == trait.Name {
					group.Traits = utils.RemoveIndex(group.Traits, index)
					break
				}
			}
			group.Points -= trait.Points
			return group, nil
		}
		err := engine.AddScript(name, body)
		if err != nil {
			return group, fmt.Errorf("%s: %s", errorPrefix, err)
		}
		script = engine.GetScript(name)
	}

	tengoizedGroup, err := scripting.ConvertToTengoMap(group)
	if err != nil {
		return group, fmt.Errorf("%s: %s", errorPrefix, err)
	}
	script.Add("group", tengoizedGroup)

	tengoizedChoices := make(map[string]any)
	for _, choice := range trait.Choices {
		tengoizedChoices[choice.Name] = choice.Value
	}
	script.Add("choices", tengoizedChoices)

	result, err := script.Run()
	if err != nil {
		return group, fmt.Errorf("%s: %s", errorPrefix, err)
	}
	output_group, err := scripting.ConvertFromTengoMap[Group](result.Get("group").Map())
	if err != nil {
		return group, fmt.Errorf("%s: %s", errorPrefix, err)
	}
	updatedGroup = &output_group
	for index, groupTrait := range updatedGroup.Traits {
		if groupTrait == trait.Name {
			updatedGroup.Traits = utils.RemoveIndex(updatedGroup.Traits, index)
			break
		}
	}
	updatedGroup.Points -= trait.Points
	return updatedGroup, nil
}

func (trait Trait) OnRemoveScriptBody() string {
	var scriptBuilder strings.Builder
	for _, change := range trait.Scripting.OnRemove {
		scriptBuilder.WriteString(fmt.Sprintf("%s\n", change))
	}
	return scriptBuilder.String()
}

func RollForCaptainTrait(availableTraits []Trait) (captainsTrait Trait, err error) {
	rand.Seed(time.Now().UnixNano())
	result, _, err := dice.Roll("3d6")
	if err != nil {
		return
	}

	for _, trait := range FilterTraitsByType("Captain", availableTraits) {
		if trait.Roll == result.Int() {
			return trait, nil
		}
	}

	err = fmt.Errorf("No available captain's trait matched roll result %d: %+v", result.Int(), availableTraits)
	return
}

func FilterTraitsBySource(sourceName string, traitList []Trait) []Trait {
	filteredTraits := []Trait{}
	for _, trait := range traitList {
		if trait.Source == sourceName {
			filteredTraits = append(filteredTraits, trait)
		}
	}
	return filteredTraits
}

func FilterTraitsByType(typeName string, traitList []Trait) []Trait {
	filteredTraits := []Trait{}
	for _, trait := range traitList {
		if trait.Type == typeName {
			filteredTraits = append(filteredTraits, trait)
		}
	}
	return filteredTraits
}

func GetTraitByName(name string, traitList []Trait) Trait {
	for _, trait := range traitList {
		if trait.Name == name {
			return trait
		}
	}
	return Trait{}
}
