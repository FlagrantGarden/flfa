package persona

import (
	"fmt"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/spf13/afero"
)

// A Persona is a generic structure representing some sort of user mode or profile that changes the internal behavior
// of an application. It must have structured Data and Settings. Personas use viper to save and load their state.
type Persona[D state.Initializable[D, *D], S state.Initializable[S, *S]] struct {
	// The Name of the Persona; this is used in messaging and to determine the name of the Persona's state file.
	// Every Persona must have a name - if one is not specified, it defaults to "default" when the Persona is
	// initialized; this is used for the name of the state file on disk as "{Name}.yaml".
	//
	// Example Behavior:
	//
	// With Name "foo":
	//     "foo.yaml"
	// With Name "foo_bar.baz":
	//     "foo_bar.baz.yaml"
	Name string
	// The Kind of Persona this is. The Name field is used in log messaging/displays about the Persona while the
	// FolderName field is used to determine where to save it and look for it when loading. If not specified, it
	// defaults to state.Kind{Name: "persona", FolderName: "personas"}
	//
	// Example Behavior:
	//
	// With root folder "myApp", Kind.FolderName "users", and Name "foo":
	//     "myApp/saves/foo.yaml"
	// With root folder "myApp", Kind.FolderName "modes", and Name "interactive":
	//     "myApp/modes/interactive.yaml"
	Kind state.Kind
	// The handler for this Persona's state to enable file discovery and save/load operations.
	Handle *state.Handle
	// Personas have Data of the type you specify when instantiating them. It must be a struct of the information related
	// to this Persona specifically. This can be any struct you want as long as:
	//
	// 1. It has an Initialize() method that returns a pointer to the data structure itself. You want to make sure that
	// this method sets sensible defaults for the data if you need any - otherwise it can just return a pointer to an
	// empty instance of the data type.
	//
	// 2. It can be marshalled/unmarshalled by viper. In general, you don't need to worry about this unless you're doing
	// something very complex. No such errors or issues have been reported, but viper is what is being used further down
	// the stack to load/save the state file for an instance so if viper can't handle it the instance will error.
	Data D
	// Settings is a struct that determins how this Persona should behave in the application. The Settings field is always
	// determined by the application developer and can be any struct so long as it has an Initialize() method which
	// returns a pointer to itself.
	// Personas have Settings of the type you specify when instantiating them and determine how this Persona should behave
	// in the application. The Settings can be any struct you want as long as:
	//
	// 1. It has an Initialize() method that returns a pointer to the settings struct itself. You want to make sure that
	// this method sets sensible defaults for the data if you need any - otherwise it can just return a pointer to an
	// empty instance of the data type.
	//
	// 2. It can be marshalled/unmarshalled by viper. In general, you don't need to worry about this unless you're doing
	// something very complex. No such errors or issues have been reported, but viper is what is being used further down
	// the stack to load/save the state file for an instance so if viper can't handle it the instance will error.
	Settings S
}

// FolderPath takes the root folder where Personas are expected to be found and appends the file-path-safe foldername
// as specified by this Persona's Kind; if no Kind is specified, the default Kind is used.
func (persona *Persona[D, S]) FolderPath(rootFolderPath string) string {
	if persona.Kind == (state.Kind{}) {
		persona.Kind = *GetDefaultKind()
	}

	return filepath.Join(rootFolderPath, utils.ValidFileName(persona.Kind.FolderName))
}

// Initialize handles the first load/set for a Persona. If the Persona's handle already has a viper, Initialize returns
// immediately. If the Persona does not have a viper, one is initialized. If the Persona's file already exists,
// Initialize will try to load it and return. If the file does not exist, Initialize will initialize the Persona's
// Settings and Data and then write the defaults to disk.
func (persona *Persona[D, S]) Initialize(name string, rootFolderPath string, afs *afero.Afero) error {
	if persona.Name == "" {
		if name != "" {
			persona.Name = name
		} else {
			persona.Name = "default"
		}
	}
	if persona.Kind == (state.Kind{}) {
		persona.Kind = *GetDefaultKind()
	}

	if persona.Handle == nil {
		persona.Handle = &state.Handle{}
	}

	handleFolderPath := persona.FolderPath(rootFolderPath)
	err := persona.Handle.Initialize(persona.Name, handleFolderPath, afs)
	if err != nil {
		return fmt.Errorf("unable to initialize %s '%s'; %s", persona.Kind.Name, persona.Name, err)
	}

	alreadyExists, err := afs.Exists(persona.Handle.FilePath)
	if err != nil {
		return fmt.Errorf("unable to initialize %s '%s': %s", persona.Kind.Name, persona.Name, err)
	}

	if alreadyExists {
		err = persona.Load(afs)
		if err != nil {
			return fmt.Errorf("unable to initialize %s: %s", persona.Kind.Name, err)
		}
		return nil
	}

	persona.Data = *persona.Data.Initialize()
	persona.Settings = *persona.Settings.Initialize()

	err = persona.Save(afs)
	if err != nil {
		return fmt.Errorf("unable to initialize %s: %s", persona.Kind.Name, err)
	}

	return nil
}

// Load attempts to read a Persona's Settings and Data from its saved state, updating the Persona if able/needed. It
// expects that the Persona's state Handle has already been initialized. If successful, it updates the Handle's
// FileInfo. If the operation fails, it returns an error.
func (persona *Persona[D, S]) Load(afs *afero.Afero) error {
	updatedFileInfo, err := persona.Handle.Load(persona, afs)
	if err != nil {
		return fmt.Errorf("unable to load %s '%s' from '%s': %s", persona.Kind.Name, persona.Name, persona.Handle.FilePath, err)
	}

	persona.Handle.FileInfo = updatedFileInfo

	return nil
}

// Save attempts to write the current in-memory representation of a Persona's Settings and Data to its file, overwriting
// any conflicting values. It expects that the Persona's state Handle has already been initialized. If successful, it
// updates the Handle's FileInfo. If unsuccessful, it returns an error.
func (persona *Persona[D, S]) Save(afs *afero.Afero) error {
	err := persona.Handle.SetStruct(persona.Settings, "settings")
	if err != nil {
		return fmt.Errorf("unable to save %s '%s' to '%s': %s", persona.Kind.Name, persona.Name, persona.Handle.FilePath, err)
	}

	err = persona.Handle.SetStruct(persona.Data, "data")
	if err != nil {
		return fmt.Errorf("unable to save %s '%s' to '%s': %s", persona.Kind.Name, persona.Name, persona.Handle.FilePath, err)
	}

	updatedFileInfo, err := persona.Handle.Save(afs)
	if err != nil {
		return fmt.Errorf("unable to save %s '%s' to '%s': %s", persona.Kind.Name, persona.Name, persona.Handle.FilePath, err)
	}

	persona.Handle.FileInfo = updatedFileInfo

	return nil
}

// DiscoverPersonas is a convenience function to look for Personas in a folder and initialize them all. The kind tells
// the function where to expect to find the personas. Because the Kind is a pointer, you can pass nil if you want; this
// will default to using the default "persona" kind, which is saved into a "personas" subfolder. The specified root
// folder path must be where you expect to find the subfolder specified by the kind parameter; this subfolder should
// contain the Persona state files.
//
// DiscoverPersonas will stop processing as soon as an error is encountered, whether that's during the initial discovery
// of the Personas or the initialization of one of them; currently there is no handling to ignore a malformed state file
// for a Persona and load the others.
//
// There is currently no auto-handling of instances associated with a Persona; each Tympan application using the
// Persona-associated instance model will need to implement their own handling.
//
// Examples:
//
// When looking in the user data folder for your app:
//
//     kind := &persona.Kind{Name: "user", FolderName: "users"}
//     persona.DiscoverPersonas[MyUserData, MyUserSettings](kind, "~/.myApp/data", afs)
//
// The function will look in "~/.myApp/data/users" for yaml files and attempt to initialize them as Personas where
// the Data field is of the `MyUserData` type and the Settings field is of the `MyUserSettings` type.
//
//     persona.DiscoverPersonas[MyUserData, MyUserSettings](nil, "~/.myApp/data", afs)
//
// The function will look in "~/.myApp/data/personas" for yaml files and attempt to initialize them as Personas where
// the Data field is of the `MyUserData` type and the Settings field is of the `MyUserSettings` type.
func DiscoverPersonas[D state.Initializable[D, *D], S state.Initializable[S, *S]](kind *state.Kind, dataFolderPath string, afs *afero.Afero) (personas []*Persona[D, S], err error) {
	if kind == nil {
		kind = GetDefaultKind()
	}

	discoveredPersonaNames, err := state.Discover(utils.ValidFileName(kind.FolderName), dataFolderPath, afs)
	if err != nil {
		return
	}

	for _, name := range discoveredPersonaNames {
		persona := &Persona[D, S]{Kind: *kind}
		err = persona.Initialize(name, dataFolderPath, afs)
		if err != nil {
			return
		}

		personas = append(personas, persona)
	}

	return
}

// GetPersona is a convenience function to find, load, and return a single Persona. You must pass the name of the
// Persona you're looking for, the Kind of Persona you're looking for, the root folder path to where you expect to find
// Personas, and an Afero file system. If you specify the Kind as nil, GetPersona will use the default Kind. GetPersona
// creates a Persona object with the specified Name and Kind and then initializes the handle to that Persona before
// attempting to Load it, returning the Persona with values populated from disk. If the load step fails, it returns a
// zero representation of the Persona and the error.
func GetPersona[D state.Initializable[D, *D], S state.Initializable[S, *S]](name string, kind *state.Kind, rootFolderPath string, afs *afero.Afero) (persona *Persona[D, S], err error) {
	if kind == nil {
		kind = GetDefaultKind()
	}

	persona = &Persona[D, S]{Name: name, Kind: *kind}
	personaFolderPath := persona.FolderPath(rootFolderPath)
	persona.Handle = &state.Handle{}
	persona.Handle.Initialize(name, personaFolderPath, afs)

	err = persona.Load(afs)
	if err != nil {
		return &Persona[D, S]{}, err
	}
	return
}

// Returns a pointer to the default Kind, which has a Name of "persona" and a FolderName of "personas"
func GetDefaultKind() *state.Kind {
	return &state.Kind{
		Name:       "persona",
		FolderName: "personas",
	}
}
