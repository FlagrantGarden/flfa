package flfa

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type BaseProfile struct {
	Source           string
	Type             string
	Category         string
	Melee            Melee
	Move             Move
	Missile          Missile
	FightingStrength FightingStrength
	Resolve          int
	Toughness        int
	Traits           []string
	Points           int
}

type BaseProfileData struct {
	Profiles []BaseProfile `mapstructure:"profiles"`
}

type BaseProfiler interface {
	Name() string
}

type Melee struct {
	Activation     int
	ToHitAttacking int
	ToHitDefending int
}

type Move struct {
	Activation int
	Distance   int
}

type Missile struct {
	Activation int
	ToHit      int
	Range      int
}

type FightingStrength struct {
	Current int
	Maximum int
}

func (ffapi *Api) ReadAndParseProfileData(dataFilePath string) ([]BaseProfile, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for Trait Data File '%s'", dataFilePath)
	}

	file, err := ffapi.AFS.ReadFile(dataFilePath)
	if err != nil {
		log.Error().Msgf("unable to read Trait Data File '%s'", dataFilePath)
	}

	// Determine source of profile:
	profileSource := filepath.Base(filepath.Dir(dataFilePath))

	var profileData BaseProfileData

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		return []BaseProfile{}, fmt.Errorf("unable to read BaseProfile Data File '%s': %s", dataFilePath, err.Error())
	}
	err = viper.Unmarshal(&profileData)

	if err != nil {
		return []BaseProfile{}, fmt.Errorf("unable to parse BaseProfile Data File '%s' %s", dataFilePath, err)
	}

	var profiles []BaseProfile

	for _, trait := range profileData.Profiles {
		trait.Source = profileSource
		profiles = append(profiles, trait)
	}

	return profiles, nil
}

func (ffapi *Api) CacheBaseProfiles(modulePath string) error {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for module at '%s'", modulePath)
	}

	moduleProfilesPath := filepath.Join(modulePath, "Profiles.yaml")
	log.Trace().Msgf("Loading profiles from %s", moduleProfilesPath)

	profiles, err := ffapi.ReadAndParseProfileData(moduleProfilesPath)
	if err != nil {
		return err
	}

	ffapi.CachedProfiles = append(ffapi.CachedProfiles, profiles...)

	return nil
}

func (bp *BaseProfile) Name() string {
	return fmt.Sprintf("%s %s", bp.Type, bp.Category)
}

func GetBaseProfile(name string, profileList []BaseProfile) (BaseProfile, error) {
	log.Trace().Msgf("searching for profile '%s'", name)
	for _, profile := range profileList {
		if profile.Name() == name {
			return profile, nil
		}
	}
	return BaseProfile{}, fmt.Errorf("no base profile found that matches name '%s'", name)
}

func GetBaseProfileTraits(profile BaseProfile, traitList []Trait) []Trait {
	traits := []Trait{}
	for _, traitName := range profile.Traits {
		trait := GetTraitByName(traitName, traitList)
		if (trait != Trait{}) {
			traits = append(traits, trait)
		}
	}
	return traits
}
