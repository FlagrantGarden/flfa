package flfa

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Spell struct {
	Source   string
	Name     string
	Check    int
	Range    int
	Target   string
	Duration string
	Effect   string
}

type SpellData struct {
	Spells []Spell `mapstructure:"spells"`
}

func (ffapi *Api) ReadAndParseSpellData(dataFilePath string) ([]Spell, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for Spell Data File '%s'", dataFilePath)
	}

	file, err := ffapi.AFS.ReadFile(dataFilePath)
	if err != nil {
		log.Error().Msgf("unable to read Spell Data File '%s'", dataFilePath)
	}

	// Determine source of spell:
	spellSource := filepath.Base(filepath.Dir(dataFilePath))

	var spellData SpellData

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		return []Spell{}, fmt.Errorf("unable to read Spell Data File '%s': %s", dataFilePath, err.Error())
	}
	err = viper.Unmarshal(&spellData)

	if err != nil {
		return []Spell{}, fmt.Errorf("unable to parse Spell Data File '%s' %s", dataFilePath, err)
	}

	var spells []Spell

	for _, spell := range spellData.Spells {
		spell.Source = spellSource
		spells = append(spells, spell)
	}

	return spells, nil
}

func (ffapi *Api) CacheModuleSpells(modulePath string) error {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for module at '%s'", modulePath)
	}

	moduleProfilesPath := filepath.Join(modulePath, "Spells.yaml")
	log.Trace().Msgf("Loading spells from %s", moduleProfilesPath)

	spells, err := ffapi.ReadAndParseSpellData(moduleProfilesPath)
	if err != nil {
		return err
	}

	ffapi.CachedSpells = append(ffapi.CachedSpells, spells...)

	return nil
}
