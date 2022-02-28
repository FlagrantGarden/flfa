package flfa

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Company struct {
	Name        string
	Description string
	Groups      []Group
	Source      string
}

type CompanyData struct {
	Companies []Company `mapstructure:"companies"`
}

func (ffapi *Api) ReadAndParseCompanyData(dataFilePath string) ([]Company, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for Company Data File '%s'", dataFilePath)
	}

	file, err := ffapi.AFS.ReadFile(dataFilePath)
	if err != nil {
		log.Error().Msgf("unable to read Company Data File '%s'", dataFilePath)
	}

	// Determine source of company:
	companySource := filepath.Base(filepath.Dir(dataFilePath))

	var companyData CompanyData

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		return []Company{}, fmt.Errorf("unable to read Company Data File '%s': %s", dataFilePath, err.Error())
	}
	err = viper.Unmarshal(&companyData)

	if err != nil {
		return []Company{}, fmt.Errorf("unable to parse Company Data File '%s' %s", dataFilePath, err)
	}

	var companies []Company

	for _, company := range companyData.Companies {
		groups := []Group{}
		for index, groupData := range company.Groups {
			group, err := ffapi.NewGroup(groupData.Name, groupData.BaseProfileName)
			if err != nil {
				return []Company{}, err
			}
			if index == 0 {
				log.Trace().Msgf("setting Group '%s' as Captain", group.Name)
				if group.Captain.Name == "" {
					group.MakeCaptain("")
				} else {
					group.Captain = groupData.Captain
				}
			}
			group.Traits = append(group.Traits, groupData.Traits...)
			log.Trace().Msgf("adding Group '%s' to Company '%s'", group.Name, company.Name)
			groups = append(groups, group)
		}
		company.Source = companySource
		company.Groups = groups
		companies = append(companies, company)
	}

	return companies, nil
}

func (ffapi *Api) CacheModuleCompanies(modulePath string) error {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for module at '%s'", modulePath)
	}

	moduleProfilesPath := filepath.Join(modulePath, "Companies.yaml")
	log.Trace().Msgf("Loading companies from %s", moduleProfilesPath)

	companies, err := ffapi.ReadAndParseCompanyData(moduleProfilesPath)
	if err != nil {
		return err
	}

	ffapi.CachedCompanies = append(ffapi.CachedCompanies, companies...)

	return nil
}
