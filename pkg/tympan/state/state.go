package state

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// A state Handle is used to save and load a stateful object to/from disk. It consists of the path to the state file,
// the latest information about that file, and a Viper instance to handle reading/writing/unmarshalling the file data.
type Handle struct {
	FilePath string
	FileInfo os.FileInfo
	Viper    *viper.Viper
}

// The Kind of a stateful struct is metadata about the struct and is specified by the application using it. The Name
// field is used in log messaging/displays about the struct while the FolderName field is used to determine where to
// save the struct or look for it when loading. If Kind is not specified, the appropriate defaults should be set by
// the implementation of the stateful struct's Initialize() method.
type Kind struct {
	Name       string
	FolderName string
}

// Stateful structs can be initialized, loaded, and saved. All three methods update the struct itself, Initialize and
// Save may write to disk while Load will only ever read from disk.
type Stateful interface {
	// A stateful struct must be initializeable. It is expected that initialize should set the default values for the
	// struct if they are not set, call Initialize on their state handle, load the state from disk if the file already
	// exists or write the default state to disk if it does not. If initalization fails, the method should return an error
	Initialize(name string, rootFolderPath string, afs *afero.Afero) error
	Load(afs *afero.Afero) error
	Save(afs *afero.Afero) error
}

// Initializable structs must be able to load their own defaults. Initializable structs are often included in Stateful
// structs. When specified in a new type of stateful struct, you must pass the type twice: once as itself, once as a
// pointer to itself. For Example:
//
//     type Instance[D state.Initializable[D, *D]] struct {}
type Initializable[I any, T any] interface {
	// The Initializable struct should return a pointer to a valid instance of itself with defaults set.
	Initialize() T
}

// MetaConfig structs are used when parsing tags on configuration structs; they help turn a mapstructure tag into the
// name of a viper configuration key and change the behavior of a configuration item via the tympanconfig directive;
// right now the only supported directive is `tympanconfig:"ignore"` which ensures a struct key is not written to the
// configuration.
type MetaConfig struct {
	ConfigKey string
	Ignore    bool
}

// ParseStructTags() is used to introspect on a struct which is to be saved to disk via viper; it returns the MetaConfig
// for a given struct field which SetStruct uses to determine behavior.
func ParseStructTags(tagEntry reflect.StructTag) (metaConfig MetaConfig) {
	mapstructTag, ok := tagEntry.Lookup("mapstructure")
	if ok {
		ignoreEntries := []string{"squash", "remain", "omitempty"}
		for _, entry := range strings.Split(mapstructTag, ",") {
			if !utils.Contains(ignoreEntries, entry) {
				metaConfig.ConfigKey = entry
			}
		}
	}
	tympanConfigTag, ok := tagEntry.Lookup("tympanconfig")
	if ok {
		tympanConfigDirectives := strings.Split(tympanConfigTag, ",")
		if utils.Contains(tympanConfigDirectives, "ignore") {
			metaConfig.Ignore = true
		}
	}
	return
}

// CurrentFileInfo returns information about a state file, including its name, size, permissions,
// and when it was last modified.
func CurrentFileInfo(v *viper.Viper, afs *afero.Afero) (os.FileInfo, error) {
	path := v.ConfigFileUsed()
	return afs.Stat(path)
}

// Initialize Viper makes sure the Handle has a viper and returns true if it had to create one or false if not.
func (handle *Handle) InitializeViper(afs *afero.Afero) bool {
	if handle.Viper == nil {
		handle.Viper = viper.New()
		handle.Viper.SetFs(afs)
		return true
	}
	return false
}

// Initialize makes sure a handle can be used to save and load a state file with viper, determining the file path
// if needed and creating the containing folder if it does not exist. It expects that the parent struct of the
// handle has already been initialized. It errors if any step fails.
func (handle *Handle) Initialize(name string, folderPath string, afs *afero.Afero) error {
	if handle == nil {
		return fmt.Errorf("Handle is nil; must be instantiated before it can be initialized")
	}
	initializedViper := handle.InitializeViper(afs)
	if !initializedViper {
		return nil
	}

	if handle.FilePath != "" {
		handle.Viper.SetConfigFile(handle.FilePath)
		return nil
	}

	if name == "" {
		return fmt.Errorf("unable to initialize handle without name")
	}

	if folderPath == "" {
		return fmt.Errorf("unable to initialize handle without folder path")
	}

	err := afs.MkdirAll(folderPath, 0755)
	if err != nil {
		return fmt.Errorf("unable to initialize handle: %s", err)
	}

	filename := fmt.Sprintf("%s.yaml", utils.ValidFileName(name))
	handle.FilePath = filepath.Join(folderPath, filename)
	handle.Viper.SetConfigFile(handle.FilePath)
	return nil
}

// Load overwrites the in-memory representation of a state from the data in a statefile. The Handle must be initialized
// before calling Load. It updates the FileInfo for the handle if successful. It returns an error if any step fails.
func (handle *Handle) Load(state Stateful, afs *afero.Afero) (os.FileInfo, error) {
	if handle.Viper == nil {
		return nil, fmt.Errorf("handle must be initialized prior to load")
	}

	err := handle.Viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = handle.Viper.Unmarshal(&state)
	if err != nil {
		return nil, err
	}

	return CurrentFileInfo(handle.Viper, afs)
}

// SetStruct is used to automatically map a struct's fields to viper configuration settings and is mapstructure-aware;
// if any of the fields in your struct have a defined mapstructure tag, SetStruct will figure out the right way to add
// them for you. It respects the `tympanconfig` tag, using any behavior specified by it. It handles nested structs and
// uses viper's dot notation to write them correctly nested, prepending the children with their parent path during
// recursion. It returns an error if the passed data is not a struct.
func (handle *Handle) SetStruct(data any, parent string) error {
	reflected_value := reflect.ValueOf(data)
	if reflected_value.Kind() == reflect.Pointer {
		reflected_value = reflect.Indirect(reflected_value)
	}
	reflected_type := reflected_value.Type()
	if reflected_value.Kind() != reflect.Struct {
		return fmt.Errorf("cannot set '%s' with SetStruct: expected a struct, got '%s'", reflected_type, reflected_value.Kind())
	}

	// Loop over the fields, writing any values to config
	for _, field := range reflect.VisibleFields(reflected_type) {
		meta := ParseStructTags(field.Tag)

		// don't write this struct field to the config
		if meta.Ignore {
			continue
		}

		// Get the value, move from pointer to value if needed and skip invalid (zero value, no point in writing)
		value := reflected_value.FieldByIndex(field.Index)
		if value.Kind() == reflect.Pointer {
			value = reflect.Indirect(value)
		}
		if value.Kind() == reflect.Invalid {
			continue
		}

		// determine the key name for yaml, replacing downcased field with value from mapstructure if specified and
		// prepending with the parent for dot-pathing the config item for viper
		name := field.Name
		if meta.ConfigKey != "" {
			name = meta.ConfigKey
		}
		if parent != "" {
			name = fmt.Sprintf("%s.%s", parent, name)
		}

		// recurse if a struct, otherwise set the value
		if value.Kind() == reflect.Struct {
			err := handle.SetStruct(value.Interface(), name)
			if err != nil {
				return err
			}
		} else {
			handle.Viper.Set(name, value.Interface())
		}
	}

	return nil
}

// Save writes a given state's in-memory data to a statefile, creating or overwriting it as needed. The Handle must be
// initialized before calling Save. It returns an error if any step fails.
func (handle *Handle) Save(afs *afero.Afero) (os.FileInfo, error) {
	if handle.Viper == nil {
		return nil, fmt.Errorf("handle must be initialized prior to save")
	}

	err := handle.Viper.WriteConfig()
	if err != nil {
		return nil, err
	}

	return CurrentFileInfo(handle.Viper, afs)
}

// Discover looks for .yaml files in a folder and returns the name of each file without the .yaml suffix.
// If the folder does not exist or cannot be read, Discover returns an error.
func Discover(stateFolderName string, folderPath string, afs *afero.Afero) (names []string, err error) {
	discoveryFolderPath := filepath.Join(folderPath, stateFolderName)
	discoveredEntries, err := afs.ReadDir(discoveryFolderPath)
	if err != nil {
		return
	}
	for _, entry := range discoveredEntries {
		if entry.IsDir() || (filepath.Ext(entry.Name()) != ".yaml") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		names = append(names, name)
	}
	return
}
