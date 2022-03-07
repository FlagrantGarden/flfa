package data

import (
	"github.com/rs/zerolog/log"
)

type Company struct {
	Name        string
	Description string
	Groups      []Group
	Source      string
}

func (company Company) WithSource(source string) Company {
	company.Source = source
	return company
}

func (company *Company) Initialize(availableProfiles []Profile, availableTraits []Trait) error {
	groups := []Group{}
	for index, groupData := range company.Groups {
		log.Trace().Msgf("initializing Group '%s' for Company '%s'", groupData.Name, company.Name)
		group, err := NewGroup(groupData.Name, groupData.ProfileName, availableProfiles)
		if err != nil {
			return err
		}
		if index == 0 {
			log.Trace().Msgf("initializing Group '%s' as Captain", group.Name)
			if group.Captain.Name == "" {
				group.MakeCaptain("", availableTraits)
			} else {
				group.Captain = groupData.Captain
			}
		}
		group.Traits = append(group.Traits, groupData.Traits...)
		log.Trace().Msgf("initialized Group '%s' for Company '%s'", group.Name, company.Name)
		groups = append(groups, group)
	}
	company.Groups = groups
	return nil
}
