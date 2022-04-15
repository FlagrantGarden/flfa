package module

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// A Cachable type is any struct which implements the necessary methods for Tympan to find and return their data to the
// application for further use.
type Cachable[T any] interface {
	// Cachable data from a module must be able to return itself with its Source (the module it comes from) set. The
	// functions in this module cannot set the Source directly due to limitations in the implementation of generic types
	// in go. If this is addressed in the future, this implementation can be simplified. For more information, see the
	// known limitations section of the go1.18 release notes (no anchor, you'll need to search for "known limitations"):
	// https://tip.golang.org/doc/go1.18
	WithSource(source string) T
}

// Some Module data includes a subtype in addition to a source. These types must implement both the Cachable type
// constraint and the methods below.
type CachableWithSubtype[T any] interface {
	Cachable[T]
	// Cachable data from a module with a subtype must be able to return itself with its Subtype (usually, the name of the
	// file it is stored in) set. The functions in this module cannot set the Subtype directly due to both the ambiguity
	// of the field name for the subtype and limitations in the implementation of generic types in go. If this is
	// addressed in the future, this implementation can be simplified. For more information, see the known limitations
	// section of the go1.18 release notes (no anchor, you'll need to search for "known limitations"):
	// https://tip.golang.org/doc/go1.18
	WithSubtype(subtype string) T
}

// ReadAndParseData must be told what data type it is looking for, given the path to the file to read, and an Afero
// file system to use. It expects that the data is stored in a slice under the "entries" key in a yaml file. It will
// look for the absolute path to the file, determine the Source for the data (as its parent folder), use viper to read
// and unmarshal the data, and then return all entries with their Source set. If any step fails, it will return an
// empty slice of the specified data type and the error.
//
// You need not call ReadAndParseData directly when interacting with a module; you can instead use the more convenient
// GetDataByFile function, which only needs the module path, the type name of the data file as it is stored on disk,
// and an Afero file system.
//
// This function is used by both GetDataByFile and GetDataByFolder.
func ReadAndParseData[T Cachable[T]](dataFilePath string, afs *afero.Afero) ([]T, error) {
	dataFilePath, err := filepath.Abs(dataFilePath)
	if err != nil {
		return []T{}, fmt.Errorf("could not find absolute path for data file '%s'", dataFilePath)
	}

	// Determine source of the data:
	source := filepath.Base(filepath.Dir(dataFilePath))

	var data struct {
		Entries []T `mapstructure:"entries"`
	}

	v := viper.New()
	v.SetFs(afs)
	v.AddConfigPath(dataFilePath)
	err = v.ReadInConfig()
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

// GetDataByFile must be told what data type it is looking for and given the path to the module folder, the name of the
// data type, and an Afero file system to use. It expects that the data is stored in a slice under the "entries" key in
// a yaml file named the same as the passed data type. It will look for the absolute path to the module folder,
// determine the name of the yaml file, and then call ReadAndParseData with the passed data type and determined file
// path, returning the slice of discovered entries. If any step fails, it will return an empty slice of the specified
// data type and the error.
func GetDataByFile[T Cachable[T]](modulePath string, dataTypeName string, afs *afero.Afero) ([]T, error) {
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

// GetDataByFolder must be told what data type it is looking for and given the path to the module folder, the name of
// the data folder to look in, and an Afero file system to use. It expects that the data is stored in multiple yaml
// files whose name is the subtype for all entries in that file. It expects that each file stores the data in a slice
// under the "entries" key. It will look for the absolute path to the module folder, combine that with the data folder
// name, and then walk the data folder, calling ReadAndParseData on each yaml file it finds, setting their subtype
// before returning the combined slice of all discovered entries from every parsed data file. If any step fails, it will
// return an empty slice of the specified data type and the error.
func GetDataByFolder[T CachableWithSubtype[T]](modulePath string, dataFolderName string, afs *afero.Afero) ([]T, error) {
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

// ReadAndParseEmbeddedData must be told what data type it is looking for, given the path to the file to read, and a
// pointer to the embedded file system to use. It expects that the data is stored in a slice under the "entries" key in
// a yaml file. It will look for the absolute path to the file, determine the Source for the data (as its parent
// folder), use viper to read and unmarshal the data, and then return all entries with their Source set. If any step
// fails, it will return an empty slice of the specified data type and the error.
//
// You need not call ReadAndParseEmbeddedData directly when interacting with a module; you can instead use the more
// convenient GetEmbeddedDataByFile function, which only needs the module path, the type name of the data file as it is
// stored on disk, and an embedded file system.
//
// This function is used by both GetEmbeddedDataByFile and GetEmbeddedDataByFolder.
func ReadAndParseEmbeddedData[T Cachable[T]](dataFilePath string, efs *embed.FS) ([]T, error) {
	// Determine source of data file:
	source := filepath.Base(filepath.Dir(filepath.Dir(dataFilePath)))

	var data struct {
		Entries []T `mapstructure:"entries"`
	}

	v := viper.New()
	dataFileBytes, err := efs.ReadFile(dataFilePath)
	if err != nil {
		return []T{}, fmt.Errorf("unable to read data file '%s': %s", dataFilePath, err.Error())
	}

	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(dataFileBytes))
	if err != nil {
		return []T{}, fmt.Errorf("unable to read data file '%s': %s", dataFilePath, err.Error())
	}

	err = v.Unmarshal(&data)
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

// GetEmbeddedDataByFile must be told what data type it is looking for and given the path to the module folder, the name
// of the data type, and an embedded file system to use. It expects that the data is stored in a slice under the
// "entries" key in a yaml file named the same as the passed data type. It will determine the name of the yaml file and
// then call ReadAndParseData with the passed data type and determined file path, returning the slice of discovered
// entries. If any step fails, it will return an empty slice of the specified data type and the error.
func GetEmbeddedDataByFile[T Cachable[T]](modulePath string, dataTypeName string, efs *embed.FS) ([]T, error) {
	dataFileName := fmt.Sprintf("%s.yaml", dataTypeName)
	// Can't use filepath.Join - on windows it uses a '\' which fails; *must* be '/'
	moduleDataFilePath := strings.Join([]string{modulePath, dataFileName}, "/")
	log.Trace().Msgf("Loading data from %s", moduleDataFilePath)

	entries, err := ReadAndParseEmbeddedData[T](moduleDataFilePath, efs)
	if err != nil {
		return []T{}, err
	}

	return entries, nil
}

// GetEmbeddedDataByFolder must be told what data type it is looking for and given the path to the module folder, the
// name of the data folder to look in, and an embedded file system to use. It expects that the data is stored in
// multiple yaml files whose name is the subtype for all entries in that file. It expects that each file stores the data
// in a slice under the "entries" key. It will join the module folder path with the data folder name and then walk the
// data folder, calling ReadAndParseEmbeddedData on each yaml file it finds, setting their subtype before returning the
// combined slice of all discovered entries from every parsed data file. If any step fails, it will return an empty
// slice of the specified data type and the error.
func GetEmbeddedDataByFolder[T CachableWithSubtype[T]](modulePath string, dataFolderName string, efs *embed.FS) ([]T, error) {
	// Can't use filepath.Join - on windows it uses a '\' which fails; *must* be '/'
	moduleFolderPath := strings.Join([]string{modulePath, dataFolderName}, "/")
	log.Trace().Msgf("Loading %s from %s", dataFolderName, moduleFolderPath)

	var returnEntries []T

	// find all entries in the module
	err := fs.WalkDir(efs, moduleFolderPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Logger.Trace().Msgf("Walking %s", path)
		isDataFile, _ := filepath.Match("*.yaml", filepath.Base(path))
		if isDataFile {
			entries, err := ReadAndParseEmbeddedData[T](path, efs)
			if err != nil {
				return err
			}
			subtype := strings.Split(filepath.Base(path), ".")[0]

			for _, entry := range entries {
				entry = entry.WithSubtype(subtype)
				returnEntries = append(returnEntries, entry)
			}
		}
		return nil
	})
	if err != nil {
		return []T{}, err
	}
	return returnEntries, nil
}
