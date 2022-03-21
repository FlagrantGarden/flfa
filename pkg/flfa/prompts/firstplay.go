package prompts

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/skirmish"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/user"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/instance"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/rs/zerolog/log"
)

func WantCustomUser() *confirmation.Confirmation {
	message := "It looks like this is your first time playing Flagrant Factions! Do you want to create a user persona now? If not, you can just use the default one."
	return confirmation.New(message, confirmation.No)
}

func GetUserName() *textinput.TextInput {
	prompt := textinput.New("What name do you want to be called by?")
	prompt.Placeholder = "Name cannot be empty"
	return prompt
}

func SelectCompany(companies []data.Company) *selection.Selection {
	var companyNames []string
	companyNames = append(companyNames, "Create a new company")
	for _, company := range companies {
		companyNames = append(companyNames, company.Name)
	}

	prompt := selection.New("Which company would you like to play as?", selection.Choices(companyNames))
	prompt.PageSize = 5

	return prompt
}

func NameCompany() *textinput.TextInput {
	prompt := textinput.New("What is your company called?")
	prompt.Placeholder = "Name cannot be empty"
	return prompt
}

func DescribeCompany() *textinput.TextInput {
	prompt := textinput.New("How would you describe your company?")
	prompt.Placeholder = "Description cannot be empty"
	return prompt
}

func StartNewSkirmish() *confirmation.Confirmation {
	message := "Would you like to start a skirmish now? If not, your info will be saved and you'll start a skirmish the next time you play."
	return confirmation.New(message, confirmation.Yes)
}

func SelectSkirmish() *selection.Selection {
	// hard code to empty void for now
	return selection.New("What scenario do you want to play?", selection.Choices([]string{"Empty Void"}))
}

func WhatNext(skirmish *instance.Instance[skirmish.Skirmish]) *selection.Selection {
	options := []string{"Check on Something", "Attack", "Move", "Cast", "Shoot", "Save & Quit"}
	return selection.New("What do you want to do now?", selection.Choices(options))
}

func FirstPlay(ffapi *flfa.Api) error {
	wantCustomPersona, err := WantCustomUser().RunPrompt()
	if err != nil {
		return err
	}

	userKind := user.Kind()
	activeUser := &persona.Persona[user.Data, user.Settings]{Kind: *userKind}
	if wantCustomPersona {
		name, err := GetUserName().RunPrompt()
		if err != nil {
			return err
		}
		log.Logger.Trace().Msgf("user name: %s", name)

		err = activeUser.Initialize(name, ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
		if err != nil {
			return err
		}
	} else {
		err = activeUser.Initialize("", ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
		if err != nil {
			return err
		}
	}

	log.Logger.Trace().Msgf("active user: %+v", activeUser)
	ffapi.Tympan.Configuration.ActiveUserPersona = activeUser.Name
	err = ffapi.Tympan.SaveConfig()
	if err != nil {
		return fmt.Errorf("error saving updated config with active user persona: %s", err)
	}

	companyChoice, err := SelectCompany(ffapi.Cache.Companies).RunPrompt()
	if err != nil {
		return err
	}
	if companyChoice.String == "Create a new company" {
		companyName, err := NameCompany().RunPrompt()
		if err != nil {
			return err
		}

		companyDescription, err := DescribeCompany().RunPrompt()
		if err != nil {
			return err
		}

		var groups []data.Group
		group, err := NewGroup(ffapi.Cache.Profiles, ffapi.Cache.Traits)
		if err != nil {
			return err
		}
		groups = append(groups, group)

		activeUser.Data.Companies = append(activeUser.Data.Companies, data.Company{
			Name:        companyName,
			Description: companyDescription,
			Groups:      groups,
			Source:      "custom",
		})
	} else {
		company, err := data.GetCompany(companyChoice.String, ffapi.Cache.Companies)
		if err != nil {
			return err
		}
		activeUser.Data.Companies = append(activeUser.Data.Companies, company)
	}

	err = activeUser.Save(ffapi.Tympan.AFS)
	if err != nil {
		return err
	}

	startSkirmish, err := StartNewSkirmish().RunPrompt()
	if err != nil {
		return err
	}

	if startSkirmish {
		skirmishName, err := SelectSkirmish().RunPrompt()
		if err != nil {
			return err
		}

		activeUser.Settings.ActiveSkirmish = skirmishName.String
		activeUser.Settings.Skirmishes = append(activeUser.Settings.Skirmishes, user.Skirmish{
			Name: skirmishName.String,
			Configuration: user.SkirmishConfiguration{
				Autosave: true,
			},
		})
		err = activeUser.Save(ffapi.Tympan.AFS)
		if err != nil {
			return err
		}

		skirmishInstance := instance.Instance[skirmish.Skirmish]{
			Kind: *skirmish.Kind(),
			Persona: instance.Persona{
				Name: activeUser.Name,
				Kind: activeUser.Kind,
			},
			Data: skirmish.Skirmish{
				Scenario:  skirmishName.String,
				Attackers: []string{activeUser.Name},
				Defenders: []string{"Computer"},
			},
		}
		err = skirmishInstance.Initialize(skirmishName.String, ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
		if err != nil {
			return err
		}

		WhatNext(&skirmishInstance).RunPrompt()
	}
	return nil
}
