package prompts

// func StartFirstSkirmish() *confirmation.Confirmation {
// 	message := "Would you like to start a skirmish now? If not, your info will be saved and you'll start a skirmish the next time you play."
// 	return confirmation.New(message, confirmation.Yes)
// }

// func FirstPlay(ffapi *flfa.Api) error {
// 	wantCustomPersona, err := prompts.WantCustomUser().RunPrompt()
// 	if err != nil {
// 		return err
// 	}

// 	userKind := user.Kind()
// 	activeUser := &persona.Persona[user.Data, user.Settings]{Kind: *userKind}
// 	if wantCustomPersona {
// 		name, err := GetUserName().RunPrompt()
// 		if err != nil {
// 			return err
// 		}
// 		log.Logger.Trace().Msgf("user name: %s", name)

// 		err = activeUser.Initialize(name, ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		err = activeUser.Initialize("", ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	log.Logger.Trace().Msgf("active user: %+v", activeUser)
// 	ffapi.Tympan.Configuration.ActiveUserPersona = activeUser.Name
// 	err = ffapi.Tympan.SaveConfig()
// 	if err != nil {
// 		return fmt.Errorf("error saving updated config with active user persona: %s", err)
// 	}

// 	companyChoice, err := SelectCompany(ffapi.Cache.Companies).RunPrompt()
// 	if err != nil {
// 		return err
// 	}
// 	if companyChoice.String == "Create a new company" {
// 		companyName, err := NameCompany().RunPrompt()
// 		if err != nil {
// 			return err
// 		}

// 		companyDescription, err := DescribeCompany().RunPrompt()
// 		if err != nil {
// 			return err
// 		}

// 		var groups []data.Group
// 		group, err := NewGroup(ffapi.Cache.Profiles, ffapi.Cache.Traits)
// 		if err != nil {
// 			return err
// 		}
// 		groups = append(groups, group)

// 		activeUser.Data.Companies = append(activeUser.Data.Companies, data.Company{
// 			Name:        companyName,
// 			Description: companyDescription,
// 			Groups:      groups,
// 			Source:      "custom",
// 		})
// 	} else {
// 		company, err := data.GetCompany(companyChoice.String, ffapi.Cache.Companies)
// 		if err != nil {
// 			return err
// 		}
// 		activeUser.Data.Companies = append(activeUser.Data.Companies, company)
// 	}

// 	err = activeUser.Save(ffapi.Tympan.AFS)
// 	if err != nil {
// 		return err
// 	}

// 	startSkirmish, err := StartFirstSkirmish().RunPrompt()
// 	if err != nil {
// 		return err
// 	}

// 	if startSkirmish {
// 		skirmishName, err := SelectLocation().RunPrompt()
// 		if err != nil {
// 			return err
// 		}

// 		activeUser.Settings.ActiveSkirmish = skirmishName.String
// 		activeUser.Settings.Skirmishes = append(activeUser.Settings.Skirmishes, user.Skirmish{
// 			Name: skirmishName.String,
// 			Configuration: user.SkirmishConfiguration{
// 				Autosave: true,
// 			},
// 		})
// 		err = activeUser.Save(ffapi.Tympan.AFS)
// 		if err != nil {
// 			return err
// 		}

// 		skirmishInstance := instance.Instance[skirmish.Skirmish]{
// 			Kind: *skirmish.Kind(),
// 			Persona: instance.Persona{
// 				Name: activeUser.Name,
// 				Kind: activeUser.Kind,
// 			},
// 			Data: skirmish.Skirmish{
// 				Scenario:  skirmishName.String,
// 				Attackers: []string{activeUser.Name},
// 				Defenders: []string{"Computer"},
// 			},
// 		}
// 		err = skirmishInstance.Initialize(skirmishName.String, ffapi.Tympan.Configuration.FolderPaths.Cache, ffapi.Tympan.AFS)
// 		if err != nil {
// 			return err
// 		}

// 		WhatNext(&skirmishInstance).RunPrompt()
// 	}
// 	return nil
// }
