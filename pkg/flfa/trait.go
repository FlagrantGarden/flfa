package flfa

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Trait struct {
	Name   string
	Type   string
	Source string
	Roll   int
	Effect string
	Points int
}
type TraitData struct {
	Traits []Trait `mapstructure:"traits"`
}

func (ffapi *Api) ReadAndParseTraitData(dataFilePath string) ([]Trait, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for Trait Data File '%s'", dataFilePath)
	}

	file, err := ffapi.AFS.ReadFile(dataFilePath)
	if err != nil {
		log.Error().Msgf("unable to read Trait Data File '%s'", dataFilePath)
	}

	// Determine type of trait:
	traitType := strings.Split(filepath.Base(dataFilePath), ".")[0]
	// Determine source of trait:
	traitSource := filepath.Base(filepath.Dir(filepath.Dir(dataFilePath)))

	var traitData TraitData

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		return []Trait{}, fmt.Errorf("unable to read Trait Data File '%s': %s", dataFilePath, err.Error())
	}
	err = viper.Unmarshal(&traitData)

	if err != nil {
		return []Trait{}, fmt.Errorf("unable to parse Trait Data File '%s' %s", dataFilePath, err)
	}

	var traits []Trait

	for _, trait := range traitData.Traits {
		trait.Type = traitType
		trait.Source = traitSource
		traits = append(traits, trait)
	}

	return traits, nil
}

func (ffapi *Api) CacheModuleTraits(modulePath string) error {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for module at '%s'", modulePath)
	}

	moduleTraitsPath := filepath.Join(modulePath, "traits")
	log.Trace().Msgf("Loading traits from %s", moduleTraitsPath)

	// find all traits in the module
	err = ffapi.AFS.Walk(moduleTraitsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Trace().Msgf("Walking %s", path)
		isDataFile, _ := filepath.Match("*.yaml", filepath.Base(path))
		if isDataFile {
			traits, err := ffapi.ReadAndParseTraitData(path)
			if err != nil {
				return err
			}

			ffapi.CachedTraits = append(ffapi.CachedTraits, traits...)
		}
		return nil
	})
	return err
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
