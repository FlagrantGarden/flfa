package tympan

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Cachable[T any] interface {
	WithSource(source string) T
}

type CachableWithSubtype[T any] interface {
	Cachable[T]
	WithSubtype(subtype string) T
}

func ReadAndParseData[T Cachable[T]](dataFilePath string, afs *afero.Afero) ([]T, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for data file '%s'", dataFilePath)
	}

	file, err := afs.ReadFile(dataFilePath)
	if err != nil {
		log.Error().Msgf("unable to read data file '%s'", dataFilePath)
	}

	// Determine source of profile:
	source := filepath.Base(filepath.Dir(dataFilePath))

	var data struct {
		Entries []T `mapstructure:"entries"`
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		return []T{}, fmt.Errorf("unable to read data file '%s': %s", dataFilePath, err.Error())
	}
	err = viper.Unmarshal(&data)

	if err != nil {
		return []T{}, fmt.Errorf("unable to parse data file '%s' %s", dataFilePath, err)
	}

	var entries []T

	for _, entry := range data.Entries {
		entry = entry.WithSource(source)
		entries = append(entries, entry)
	}

	return entries, nil
}

func GetModuleDataByFile[T Cachable[T]](modulePath string, dataTypeName string, afs *afero.Afero) ([]T, error) {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		log.Error().Msgf("could not find absolute path for module at '%s'", modulePath)
	}

	dataFileName := fmt.Sprintf("%s.yaml", dataTypeName)
	moduleDataFilePath := filepath.Join(modulePath, dataFileName)
	log.Trace().Msgf("Loading data from %s", moduleDataFilePath)

	entries, err := ReadAndParseData[T](moduleDataFilePath, afs)
	if err != nil {
		return []T{}, err
	}

	return entries, nil
}

func GetModuleDataByFolder[T CachableWithSubtype[T]](modulePath string, dataFolderName string, afs *afero.Afero) ([]T, error) {
	modulePath, err := filepath.Abs(modulePath)
	if err != nil {
		return []T{}, fmt.Errorf("could not find absolute path for module at '%s'", modulePath)
	}

	moduleFolderPath := filepath.Join(modulePath, dataFolderName)
	log.Trace().Msgf("Loading %s from %s", dataFolderName, moduleFolderPath)

	var returnEntries []T

	// find all entries in the module
	err = afs.Walk(moduleFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Trace().Msgf("Walking %s", path)
		isDataFile, _ := filepath.Match("*.yaml", filepath.Base(path))
		if isDataFile {
			entries, err := ReadAndParseData[T](path, afs)
			if err != nil {
				return err
			}
			subtype := strings.Split(filepath.Base(path), ".")[0]

			for _, entry := range entries {
				entry = entry.WithSubtype(subtype)
				returnEntries = append(returnEntries, entry)
			}

			returnEntries = append(returnEntries, entries...)
		}
		return nil
	})
	return returnEntries, nil
}
