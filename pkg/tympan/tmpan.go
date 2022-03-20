package tympan

import (
	"io/fs"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/spf13/afero"
)

// The main entrypoint to using Tympan in an application, this struct handles the file system interactions, the
// configuration of the application itself, and stores the metadata for the application.
type Tympan[Config Configurable] struct {
	// Handles calls to read to/write from the filesystem. Using Afero enables easier testing and numerous alternate file
	// systems without requiring additional coding, including an in-memory FS, Google Cloud Storage, and more.
	AFS *afero.Afero
	// Configuration holds the merged configuration items from Tympan's SharedConfig and an applications own Configuration
	// struct. For more information, see the Configurable interface.
	Configuration Config
	// Metadata holds information about the application itself as set by the developer.
	Metadata Metadata
	// ConfigHandler utilizes state.Handle to read and write the configuration file for the application.
	ConfigHandler *state.Handle
}

// Metadata is set by the Tympan application developer, not the end user. These values are held for reuse in messaging,
// finding/setting defaults, and more.
type Metadata struct {
	// The short name of the application - the name of the *binary* itself without any file extensions.
	Name string
	// The name of the application when being displayed. So for the "myapp" binary, it might be "My Application"
	DisplayName string
	// The description of the application to be displayed in the help for the root command and elsewhere.
	Description string
	// Where the application's files should be stored when using the configuration or cache folders.
	FolderName string
	// The name of the application's config file; defaults to "{name}-config.yaml"
	ConfigFileName string
	// The permissions that the app should use when creating new files and folders; defaults to 0755.
	DefaultPermissions fs.FileMode
	// The URL where the application's source code can be found.
	SourceUrl string
	// The URL where the application's website (including docs) can be found.
	ProjectUrl string
}

// A reimplementation of a Tympan application must implement a small subset of functionality.
type Tympaner interface {
	InitializeConfig() error
	LoadConfig() error
	SaveConfig() error
	CachePersonas() error
}

// Returns the path to the default configuration file. A user may specify their own alternate configuration file when
// calling the application, but this returns where the app should look by default.
func (tympan *Tympan[config]) DefaultConfigFile() string {
	return filepath.Join(tympan.Configuration.GetFolderPath("configuration"), tympan.Metadata.ConfigFileName)
}
