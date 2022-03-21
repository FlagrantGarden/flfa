package tympan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

// The SharedConfig is used by all Tympan applications (those which utilize the root forme). The SharedConfig must be
// embedded in the application's own configuration struct.
type SharedConfig struct {
	// Where the application, configuration, and cache folders can be found.
	FolderPaths Folders `mapstructure:"folder_paths"`
	// What verbosity the application logging should use if not overridden by user-specified flag.
	DefaultLogLevel string `mapstructure:"default_log_level"`
	// What output format the application should use if not overridden by user-specified flag.
	DefaultOutputFormat string `mapstructure:"default_output_format"`
	// The path to the configuration file for the application. This field is included for convenience but is not written
	// to the configuration file on disk, as it represents that file for the application to find.
	ConfigurationFilePath string `tympanconfig:"ignore"`
}

// The folders where the application may place files and folders of its own.
type Folders struct {
	// This is where the application's binary is stored.
	Application string
	// This is where the application's configuration is stored by default. Other configurations may be placed
	// in this folder.
	Configuration string
	// This is where the application's state is stored by default. If using modules, plugins, or updateable documentation,
	// the files and folders for those items are stored in this folder as well.
	Cache string
}

// Reimplementing the SharedConfig requires a small subset of functions that Tympan expects to use.
type SharedConfigI interface {
	// Set default values for the shared config and potentially change the system but writing files or folders.
	InitializeSharedConfig(metadata Metadata, afs *afero.Afero) error
	// Must return the full path to a specified folder type.
	GetFolderPath(folder string) string
}

// When implementing your own configuration struct for a Tympan application, you must include a method to set default
// values for your own configuration. This is called by tympan.InitializeConfig().
//
// If you have no initialization steps or important default values, you can implement your configuration like this:
//
//      type Configuration struct {
//        tympan.SharedConfig `mapstructure:",squash" tympanconfig:"ignore"`
//        MyAppSetting        string `mapstructure:"my_app_setting"`
//      }
//
//      func (config *Configuration) Initialize() error {
//        return nil
//      }
type ApplicationConfigI interface {
	Initialize() error
}

// This interface is a type constraint because Tympan leverages go's generics. To be a valid struct which Tympan will
// know how to use for configuring an application, your struct must match the ApplicationConfigI interface and embed
// Tympan's own SharedConfig struct *or* a valid reimplementation of that struct.
//
// For example:
//      type Configuration struct {
//        tympan.SharedConfig `mapstructure:",squash" tympanconfig:"ignore"`
//        MyAppSetting        string `mapstructure:"my_app_setting"`
//      }
//
//      func (config *Configuration) Initialize() error {
//        return nil
//      }
// A few things to note about this example:
//
// 1. The SharedConfig is an embedded struct and uses two tags: `mapstructure:",squash"` and `tympanconfig:"ignore"`.
//    Both of these are necessary for your configuration to be read correctly from disk. Combined, they tell Tympan:
//    All of the values for this struct should be merged into the top-level of the struct and do not attempt to write
//    this field name to the configuration file. Both are required for an embedded struct in a state file.
//
// 2. The MyAppSetting field has a mapstructure tag of "my_app_setting". This will cause Tympan to save this field's
//    value into the config under the "my_app_setting" key instead of "myappsetting". We strongly encourage using snake-
//    cased mapstructure definitions for complex field names as viper/mapstructure will merely downcase by default,
//    which can be difficult for users to read/understand.
type Configurable interface {
	// The general configuration of a Tympan application.
	SharedConfigI
	// The configuration specific to a particular Tympan application as defined by the application developer.
	ApplicationConfigI
}

// The InitializeConfig() method handles setting default values, creating required folders, instantiating the state
// handler for the application configuration, and creating the configuration file with default values (if it does not
// already exist) or loading the configuration file's current values (if it does). It returns an error if any step of
// the process fails.
func (t *Tympan[Config]) InitializeConfig() error {
	log.Logger.Trace().Msgf("Initializing shared configuration")
	err := t.Configuration.InitializeSharedConfig(t.Metadata, t.AFS)
	if err != nil {
		return err
	}

	log.Logger.Trace().Msgf("Initializing application configuration")
	t.Configuration.Initialize()
	if err != nil {
		return err
	}

	if t.ConfigHandler == nil {
		t.ConfigHandler = &state.Handle{}
	}
	err = t.ConfigHandler.Initialize(t.Metadata.ConfigFileName, t.Configuration.GetFolderPath("configuration"), t.AFS)
	if err != nil {
		return err
	}
	log.Logger.Trace().Msgf("using configuration file: %s", t.ConfigHandler.Viper.ConfigFileUsed())

	alreadyExists, err := t.AFS.Exists(t.ConfigHandler.FilePath)
	if err != nil {
		return fmt.Errorf("unable to initialize configuration: %s", err)
	}

	if alreadyExists {
		err = t.LoadConfig()
		if err != nil {
			return fmt.Errorf("unable to initialize configuration; %s", err)
		}
		log.Logger.Trace().Msgf("loaded configuration from disk")
		return nil
	}

	err = t.SaveConfig()
	if err != nil {
		return fmt.Errorf("unable to initialize configuration; %s", err)
	}
	log.Logger.Trace().Msgf("saved configuration to disk")

	return nil
}

// The LoadConfig() method reads the Tympan application's configuration file and updates the in-memory representation of
// the configuration used in the application. It expects that the configuration has already been loaded. It returns an
// error if any step of the process fails.
func (t *Tympan[Config]) LoadConfig() error {
	if t.ConfigHandler.Viper == nil {
		return fmt.Errorf("unable to load configuration: config handler must be initialized prior to load")
	}

	err := t.ConfigHandler.Viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("unable to load configuration from '%s': %s", t.ConfigHandler.FilePath, err)
	}

	err = t.ConfigHandler.Viper.Unmarshal(&t.Configuration)
	if err != nil {
		return fmt.Errorf("unable to load configuration from '%s': %s", t.ConfigHandler.FilePath, err)
	}

	updatedFileInfo, err := state.CurrentFileInfo(t.ConfigHandler.Viper, t.AFS)
	if err != nil {
		return fmt.Errorf("unable to load configuration from '%s': %s", t.ConfigHandler.FilePath, err)
	}

	t.ConfigHandler.FileInfo = updatedFileInfo

	return nil
}

// The LoadConfig() method writes the current in-memory configuration of the Tympan application to disk. It expects that
// the configuration has already been loaded. It returns an error if any step of the process fails.
func (t *Tympan[Config]) SaveConfig() error {
	if t.ConfigHandler.Viper == nil {
		return fmt.Errorf("unable to save configuration: config handler must be initialized prior to save")
	}

	err := t.ConfigHandler.SetStruct(t.Configuration, "")
	if err != nil {
		return fmt.Errorf("unable to save configuration: %s", err)
	}

	err = t.ConfigHandler.Viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to save configuration to '%s': %s", t.ConfigHandler.FilePath, err)
	}

	updatedFileInfo, err := state.CurrentFileInfo(t.ConfigHandler.Viper, t.AFS)
	if err != nil {
		return fmt.Errorf("unable to save configuration to '%s': %s", t.ConfigHandler.FilePath, err)
	}

	t.ConfigHandler.FileInfo = updatedFileInfo

	return nil
}

// The InitializeSharedConfig() method sets the defaults for the SharedConfig struct which is always embedded in a
// Tympan application's own configuration struct. You must pass it the metadata for your Tympan application and a valid
// implementation of an Afero filesystem. This method:
//
// 1. Sets FolderPaths.Application to the parent folder containing the application's binary, returning an error if the
// location cannot be discovered or does not exist.
//
// 2. Sets FolderPaths.Configuration to the joined path of the current user's configuration directory per their
// operating system's file and folder layout standards with the application's FolderName (as defined in the Metadata) as
// the child folder. If needed, creates this folder and any missing parent folders, specifying the default permissions
// from the Metadata. It returns an error if the user configuration directory cannot be determined or an error occurs
// when ensuring the folder exists.
//
// 3. Sets FolderPaths.Cache to the joined path of the current user's cache directory per their operating system's file
// and folder layout standards with the application's FolderName (as defined in the Metadata) as the child folder. If
// needed, it creates this folder and any missing parent folders, specifying the default permissions from the Metadata.
// It returns an error if the user cache directory cannot be determined or an error occurs when ensuring the folder
// exists.
func (config *SharedConfig) InitializeSharedConfig(metadata Metadata, afs *afero.Afero) error {
	executableFilePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to set application information: %s", err)
	}
	config.FolderPaths.Application = filepath.Dir(executableFilePath)
	appFolderExists, err := afs.DirExists(config.FolderPaths.Application)
	if !appFolderExists {
		return fmt.Errorf("unable to set application information: expected folder '%s' to exist", config.FolderPaths.Application)
	} else if err != nil {
		return fmt.Errorf("unable to set application information; error accessing application folder: %s", err)
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("unable to set application information; error finding user configuration directory: %s", err)
	}
	config.FolderPaths.Configuration = filepath.Join(userConfigDir, metadata.FolderName)
	// Create config folder if needed
	err = afs.MkdirAll(config.FolderPaths.Configuration, metadata.DefaultPermissions)
	if err != nil {
		return fmt.Errorf("unable to intialize configuration folder path at '%s'; %s", config.FolderPaths.Configuration, err)
	}

	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	config.FolderPaths.Cache = filepath.Join(cacheDirectory, metadata.FolderName)
	err = afs.MkdirAll(config.FolderPaths.Cache, metadata.DefaultPermissions)
	if err != nil {
		return fmt.Errorf("unable to intialize cache folder path at '%s'; %s", config.FolderPaths.Cache, err)
	}

	return nil
}

// This method returns one of the folder paths from the shared configuration; because the Tympan library itself has to
// reference the application configuration by its type constraint it cannot use the values from the SharedConfig struct
// directly (at least not at the time of this writing and using go 1.18).
//
// Examples:
//
//     myConfiguration.GetFolderPath("configuration") // returns myConfiguration.FolderPaths.Configuration
//
//     myConfiguration.GetFolderPath("application") // returns myConfiguration.FolderPaths.Application
//
//     myConfiguration.GetFolderPath("cache") // returns myConfiguration.FolderPaths.Cache
//
//     myConfiguration.GetFolderPath("does not exist") // returns ""
func (config *SharedConfig) GetFolderPath(folder string) string {
	switch strings.ToLower(folder) {
	case "configuration":
		return config.FolderPaths.Configuration
	case "application":
		return config.FolderPaths.Application
	case "cache":
		return config.FolderPaths.Cache
	default:
		return ""
	}
}
