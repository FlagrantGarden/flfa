package prompt

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
)

type PromptContent struct {
	ErrorMessage string
	Label        string
	Result       string
}

func (p *PromptContent) GetInput() error {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(p.ErrorMessage)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     p.Label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	log.Trace().Msgf("Input: %s\n", result)
	p.Result = result

	return nil
}

func (pc *PromptContent) GetSelection(items []string) error {
	index := -1
	var result string
	var err error
	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.Label,
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	log.Trace().Msgf("Input: %s\n", result)
	pc.Result = result

	return nil
}

func (pc *PromptContent) GetConfirmation() {
	prompt := promptui.Prompt{
		Label:     pc.Label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		pc.Result = "no"
	} else {
		pc.Result = "yes"
	}

	log.Trace().Msgf("Input: %s\n", pc.Result)

}
