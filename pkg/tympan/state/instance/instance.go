package instance

import (
	"fmt"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/spf13/afero"
)

// An instance is a generic structure representing some sort of state for an application. Instances must have an
// initializable Data property. Instances can be initialized, loaded, and saved.
type Instance[D state.Initializable[D, *D]] struct {
	// Every instance must have a name - if one is not specified, it defaults to "default" when the Instance is
	// initialized; this is used for the name of the state file on disk as "{Name}.yaml".
	//
	// Example Behavior:
	//
	// With Name "foo":
	//     "foo.yaml"
	// With Name "foo_bar.baz":
	//     "foo_bar.baz.yaml"
	Name string
	// Every instance must have a Kind - if one is not specified, it defaults to "saves" when the Instance is
	// initialized; this is used to determine the folder to save the state file in as "{root folder}/{Kind}/"
	//
	// Example Behavior:
	//
	// With root folder "myApp", Kind.FolderName "saves", and Name "foo":
	//     "myApp/saves/foo.yaml"
	// With root folder "myApp", Kind.FolderName "journals", and Name "2022-01-31":
	//     "myApp/journals/2022-01-31.yaml"
	Kind state.Kind
	// Instances *may* be related to a Persona; if they are, the Persona's name is treated as a namespace for the
	// instance, using "{root folder}/personas/{Persona}/{Category}/" instead of "{root folder}/{Category}/"
	//
	// Example Behavior:
	//
	// With root folder "myApp", Kind.FolderName "saves", and Name "foo":
	//     "myapp/saves/foo.yaml"
	// With root folder "myApp", Kind.FolderName "journals", Persona.Name "foo", Persona.Kind.FolderName "authors",
	// and Name "bar":
	//     "myapp/personas/foo/saves/bar.yaml"
	Persona Persona
	// Instances have a Handle to abstract managing the loading and saving of their state and do not need to be
	// interacted with directly.
	Handle *state.Handle
	// Instances have data of the type you specify when instantiating them. This is what is written to and loaded from
	// the state file. This can be any sort of data you want as long as:
	//
	// 1. It has an Initialize() method that returns a pointer to the data structure itself. You want to make sure that
	// this method sets sensible defaults for the data if you need any - otherwise it can just return a pointer to an
	// empty instance of the data type.
	//
	// 2. It can be marshalled/unmarshalled by viper. In general, you don't need to worry about this unless you're doing
	// something very complex. No such errors or issues have been reported, but viper is what is being used further down
	// the stack to load/save the state file for an instance so if viper can't handle it the instance will error.
	Data D
}

// The Instance Persona is a minimized representation of the Persona related to an Instance, if any. It has only the
// name of the Persona and its Kinda information so that it can be located and associated.
type Persona struct {
	Name string
	Kind state.Kind
}

// FolderPath takes the root folder where Instances are expected to be found and appends the file-path-safe foldername
// as specified by this Instance's Kind; if no Kind is specified, the default Kind is used. Additionally, if this
// Instance is associated to a Persona (ie, instance.Persona is set), it will inject the associated Persona's path
// segment between the root folder and the Instance's folder.
//
// To clarify, an Instance without an associated Persona uses this formula for the FolderPath:
//     "{rootFolderPath}/{Instance.Kind.FolderName}"
// And an Instance with an associated Persona uses this formula:
//     "{rootFolderPath}/{Persona.Kind.FolderName}/{Persona.Name}/{Instance.Kind.FolderName}"
func (instance *Instance[D]) FolderPath(rootFolderPath string) string {
	if instance.Kind == (state.Kind{}) {
		instance.Kind = *GetDefaultKind()
	}

	if instance.Persona != (Persona{}) {
		return filepath.Join(rootFolderPath,
			utils.ValidFileName(instance.Persona.Kind.FolderName),
			utils.ValidFileName(instance.Persona.Name),
			utils.ValidFileName(instance.Kind.FolderName),
		)
	}

	return filepath.Join(rootFolderPath, utils.ValidFileName(instance.Kind.FolderName))
}

// Initialize handles the first load/set for an Instance. If the Instance's handle already has a viper, Initialize
// returns immediately. If the Instance does not have a viper, one is initialized. If the Instance's file already
// exists, Initialize will try to load it and return. If the file does not exist, Initialize will initialize the
// Instance's Data and then write the defaults to disk.
func (instance *Instance[D]) Initialize(name string, rootFolderPath string, afs *afero.Afero) error {
	if instance.Name == "" {
		if name != "" {
			instance.Name = name
		} else {
			instance.Name = "default"
		}
	}

	if instance.Kind == (state.Kind{}) {
		instance.Kind = *GetDefaultKind()
	}

	if instance.Handle == nil {
		instance.Handle = &state.Handle{}
	}
	handleFolderPath := instance.FolderPath(rootFolderPath)
	err := instance.Handle.Initialize(instance.Name, handleFolderPath, afs)
	if err != nil {
		return fmt.Errorf("unable to initialize %s instance '%s'; %s", instance.Kind, instance.Name, err)
	}

	alreadyExists, err := afs.Exists(instance.Handle.FilePath)
	if err != nil {
		return fmt.Errorf("unable to initialize %s instance '%s': %s", instance.Kind, instance.Name, err)
	}

	if alreadyExists {
		err = instance.Load(afs)
		return fmt.Errorf("unable to initialize %s instance: %s", instance.Kind, err)
	}

	instance.Data = *instance.Data.Initialize()
	if err != nil {
		return fmt.Errorf("unable to initialize %s instance '%s': %s", instance.Kind, instance.Name, err)
	}

	err = instance.Save(afs)
	if err != nil {
		return fmt.Errorf("unable to initialize %s instance: %s", instance.Kind, err)
	}

	return nil
}

// Load attempts to read an Instance's Data from its saved state, updating the Instance if able/needed. It expects that
// the Instance's state Handle has already been initialized. If successful, it updates the Handle's FileInfo. If the
// operation fails, it returns an error.
func (instance *Instance[D]) Load(afs *afero.Afero) error {
	updatedFileInfo, err := instance.Handle.Load(instance, afs)
	if err != nil {
		return fmt.Errorf("unable to load instance '%s' from '%s': %s", instance.Name, instance.Handle.FilePath, err)
	}

	instance.Handle.FileInfo = updatedFileInfo

	return nil
}

// Save attempts to write the current in-memory representation of an Instance's Data to its file, overwriting any
// conflicting values. It expects that the Instance's state Handle has already been initialized. If successful, it
// updates the Handle's FileInfo. If unsuccessful, it returns an error.
func (instance *Instance[D]) Save(afs *afero.Afero) error {
	err := instance.Handle.SetStruct(instance.Data, "data")
	if err != nil {
		return fmt.Errorf("unable to save %s '%s' to '%s': %s", instance.Kind.Name, instance.Name, instance.Handle.FilePath, err)
	}

	updatedFileInfo, err := instance.Handle.Save(afs)
	if err != nil {
		return fmt.Errorf("unable to save instance '%s' to '%s': %s", instance.Name, instance.Handle.FilePath, err)
	}

	instance.Handle.FileInfo = updatedFileInfo

	return nil
}

// DiscoverInstances is a convenience function to look for Instances in a folder and initialize them all. The specified
// root folder path must be where you expect to find a folder containing the state files; Instance state files *always*
// live in a folder named for their category but they may be associated to a persona and, if so, will be found in that
// persona's folder.
//
// DiscoverInstances will stop processing as soon as an error is encountered, whether that's during the initial
// discovery of the Instances or the initialization of one of them; currently there is no handling to ignore a malformed
// state file for an Instance and load the others.
//
// There is currently no auto-handling of instances associated with a Persona; each Tympan application using the
// Persona-associated instance model will need to implement their own handling.
//
// Examples:
//
// When looking in the user data folder for your app: DiscoverInstances[MySaveGame]("saves", "~/.myApp/data", afs)
//
// The function will look in "~/.myApp/data/saves" for yaml files and attempt to initialize them as Instances where the
// Data field is of the `MySaveGame` type.
//
// When looking in a specific persona: DiscoverInstances[MySaveGame]("saves", "~/.myApp/data/personas/someone", afs)
//
// The function will look in "~/.myApp/data/personas/someone/saves" for yaml files and attempt to initialize them as
// Instances where the Data field is of the `MySaveGame` type.
func DiscoverInstances[D state.Initializable[D, *D]](kind state.Kind, rootFolderPath string, afs *afero.Afero) (instances []*Instance[D], err error) {
	discoveredInstanceNames, err := state.Discover(utils.ValidFileName(kind.FolderName), rootFolderPath, afs)
	if err != nil {
		return
	}

	for _, name := range discoveredInstanceNames {
		instance := &Instance[D]{
			Name: name,
			Kind: kind,
		}

		err = instance.Initialize(name, rootFolderPath, afs)
		if err != nil {
			return
		}

		instances = append(instances, instance)
	}
	return
}

// GetInstance is a convenience function to find, load, and return a single Instance. You must pass the name of the
// Instance you're looking for, the Kind of Instance you're looking for, the instance.Persona (if any, nil if not) of
// the associated Persona for the Instance, the root folder path to where you expect to find instances, and an Afero
// file system. If you specify the Kind as nil, GetInstance will use the default Kind. If you specify the associated
// Persona as nil, the function will not look for the Instance in that Persona's folder. GetInstance creates an Instance
// object with the specified Name, Kind, and Persona and then initializes the handle to that Instance before attempting
// to Load it, returning the Instance with values populated from disk. If the load step fails, it returns a zero
// representation of the Instance and the error.
func GetInstance[D state.Initializable[D, *D]](name string, kind *state.Kind, persona *Persona, rootFolderPath string, afs *afero.Afero) (instance *Instance[D], err error) {
	if kind == nil {
		kind = GetDefaultKind()
	}

	instance = &Instance[D]{Name: name, Kind: *kind}
	if persona != nil {
		instance.Persona = *persona
	}
	instanceFolderPath := instance.FolderPath(rootFolderPath)
	instance.Handle = &state.Handle{}
	instance.Handle.Initialize(name, instanceFolderPath, afs)

	err = instance.Load(afs)
	if err != nil {
		return &Instance[D]{}, err
	}
	return
}

// Returns a pointer to the default Kind, which has a Name of "save" and a FolderName of "saves"
func GetDefaultKind() *state.Kind {
	return &state.Kind{
		Name:       "save",
		FolderName: "saves",
	}
}
