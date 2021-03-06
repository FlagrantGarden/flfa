package data

import (
	"fmt"

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
				group.PromoteToCaptain(nil, FilterTraitsBySource("core", availableTraits)...)
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

func (company Company) Points() (points int) {
	for _, group := range company.Groups {
		points += group.Points
	}

	return
}

func GetCompany(name string, companyList []Company) (Company, error) {
	log.Trace().Msgf("searching for company '%s'", name)
	for _, company := range companyList {
		if company.Name == name {
			return company, nil
		}
	}
	return Company{}, fmt.Errorf("no company found that matches name '%s'", name)
}
