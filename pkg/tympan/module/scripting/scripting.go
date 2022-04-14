package scripting

import (
	"fmt"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

// The LibraryImporter is an implementation of the tengo.ModuleGetter; it allows the engine to import dynamically
// defined libraries from modules in addition to tengo's in-the-box libraries.
type LibraryImporter struct {
	mods     tengo.ModuleGetter
	fallback func(name string) tengo.Importable
}

// The Engine is the main interface between tengo and a Tympan app and is geared towards loading scripts and libraries
// from a Tympan module.
type Engine struct {
	// The settings for the engine determine its overall behavior
	Settings EngineSettings
	// The importer is used to provide access to dynamic/module libraries and enables caching them on discovery
	Importer *LibraryImporter
	// The array of scripts enables you to cache a script for future use and introspection
	Scripts []*Script
}

// The EngineSettings configure how the tengo script engine behaves, providing some useful shorthands so you don't need
// to understand the tengo interop model in detail.
type EngineSettings struct {
	// The ScriptHeader is prepended to every script the engine caches; it ensures that any allowed/added libraries are
	// available to the script without every script needing to redeclare them.
	ScriptHeader string
	// The RandomSeed allows you to specify a seed to initialize for randomization in go instead of in your tengo scripts.
	RandomSeed int64
	// If specified, MaximumObjectAllocations limits the number of objects any one script can create.
	MaximumObjectAllocations int64
	// By default, access to tengo's OS library is forbidden. If you want to enable it, set this to true. The OS library
	// includes functions for modifying the system state, including files, folders, environment, and running arbitrary
	// processes.
	AllowOSLibrary bool
	// The list of tengo's standard libraries that the engine should cache and make available to scripts.
	StandardLibraries []string
	// The list of standalone libraries that the engine should cache and make available to scripts. These can be any tengo
	// file, so long as it exports.
	ApplicationLibraries []Library
	// The list of Tympan scripting modules that the engine should cache and make available to scripts.
	ApplicationModules []Module
	// The list of standard libraries that a script can have utilize.
	ValidStandardLibraryNames []string
}

// The Script is a wrapper around a tengo.Script object, appending the name of the script and the raw string for its
// body so that a script can be found by name and/or introspected on after compiling.
type Script struct {
	// The Name of the script. Must be unique to this instance of the Engine.
	Name string
	// The script body as a string for readability and introspection during debug.
	Body string
	// The actual tengo script object
	*tengo.Script
}

// A Tympan scripting module is a Library with zero or more submodules. This enables you to split your scripting across
// multiple files for ease of development, testing, and maintenance.
type Module struct {
	// The module library itself with an engine-unique name and the tengo script body.
	Library
	// Modules may have zero or more submodules. Submodules are made available to the module library's scope.
	Submodules []Library
}

// A Tympan scripting Library is made up of a unique name and the tengo script contents are its body.
// It is expected that a Library exports.
type Library struct {
	// The engine-unique name of the Library
	Name string
	// The tengo script that makes up the Library
	Body string
}

// Creates a new instance of an engine with no settings except that it prepopulates the list of valid standard libraries
// from tengo. Note that even though the list of valid names includes the OS library, it is disabled by default and none
// of the libraries are made available to the engine by default; this just prevents you from having to look them up.
func NewEngine() *Engine {
	return &Engine{
		Settings: EngineSettings{
			ValidStandardLibraryNames: stdlib.AllModuleNames(),
		},
	}
}

// ?
func (importer *LibraryImporter) Get(name string) tengo.Importable {
	if mod := importer.mods.Get(name); mod != nil {
		return mod
	}
	return importer.fallback(name)
}

// Returns the list of Tympan scripting libraries the engine is currently configured to be able to load
func (engine *Engine) ApplicationLibraryNames() (names []string) {
	for _, library := range engine.Settings.ApplicationLibraries {
		names = append(names, library.Name)
	}
	return names
}

// Returns the list of Tympan scripting modules the engine is currently configured to be able to load
func (engine *Engine) ApplicationModuleNames() (names []string) {
	for _, module := range engine.Settings.ApplicationModules {
		names = append(names, module.Name)
	}
	return names
}

// The InitializeLibraryImporter method is used to create and add the LibraryImporter to the engine so you don't have to
// do it manually; it is configured by default to add all of the configured standard libraries and handle importing the
// Tympan scripting libraries and modules as well. The importer is what ensures the functions, variables, etc in these
// libraries are made available to the scripts the engine will run.
func (engine *Engine) InitializeLibraryImporter() {
	// if already initialized, dont do anything
	if engine.Importer != nil {
		return
	}

	engine.Importer = &LibraryImporter{
		mods: stdlib.GetModuleMap(engine.Settings.StandardLibraries...),
		fallback: func(name string) tengo.Importable {
			// Set the source to an empty string to keep from setting things on fire
			source := ""
			// loop first over application libraries; if the name specified matches,
			// set the source to that library's body and break.
			for _, library := range engine.Settings.ApplicationLibraries {
				if library.Name == name {
					source = library.Body
					break
				}
			}

			// if the name wasn't in one of the application libraries,
			// it might be in a module or submodule
			if source == "" {
			out:
				for _, module := range engine.Settings.ApplicationModules {
					if module.Name == name {
						source = module.Body
						break
					} else {
						// if the name isn't of an application library or module, it might be of a submodule;
						// loop over the submodules for each module and check those too.
						for _, moduleLibrary := range module.Submodules {
							if moduleLibrary.Name == name {
								source = moduleLibrary.Body
								// break all the way out of the nested loop
								break out
							}
						}
					}
				}
			}

			// this should be the source of the specified library or an empty string
			return &tengo.SourceModule{Src: []byte(source)}
		},
	}
}

// AllowedStandardLibraries is a helper function for returning the list of tengo's standard libraries that the engine
// is currently configured to allow scripts to access.
func (engine *Engine) AllowedStandardLibraries() (libraries []string) {
	libraries = stdlib.AllModuleNames()
	if !engine.Settings.AllowOSLibrary {
		for index, library := range libraries {
			if library == "os" {
				libraries = utils.RemoveIndex(libraries, index)
				break
			}
		}
	}
	return libraries
}

// AddStandardLibrary adds the specified tengo standard library to the engine's cache and appends the declaration for
// using the standard library to the script header. For example, when adding the text standard library:
//
//    myengine.AddStandardLibrary("text")
//
// That call will add the text standard library to the StandardLibraries list in the engine's settings and append the
// string `text := import("text")` to the script header.
//
// Because the library is declared in the script header, any script can then use the text library without having to
// redeclare it itself.
func (engine *Engine) AddStandardLibrary(libraryName string) error {
	if utils.Contains(engine.Settings.ValidStandardLibraryNames, libraryName) {
		engine.Settings.StandardLibraries = append(engine.Settings.StandardLibraries, libraryName)
		engine.Settings.ScriptHeader += fmt.Sprintf("%s := import(\"%s\")\n", libraryName, libraryName)
		return nil
	}
	return fmt.Errorf("unable to add '%s' as standard library; must be one of: %s",
		libraryName,
		&engine.Settings.ValidStandardLibraryNames,
	)
}

// SetStandardLibraries is a helper method which replaces the existing list of standard libraries the engine is
// configured to use. It is functionally equivalent to calling AddStandardLibrary in a loop but with a little more
// safety and validation, preventing partial updates of the setting.
func (engine *Engine) SetStandardLibraries(libraryNames []string) error {
	var invalidLibraryNames []string
	for _, libraryName := range libraryNames {
		if !utils.Contains(engine.Settings.ValidStandardLibraryNames, libraryName) {
			invalidLibraryNames = append(invalidLibraryNames, libraryName)
		}
	}
	if len(invalidLibraryNames) == 1 {
		return fmt.Errorf("unable to add '%s' as standard library; must be one of: %s",
			invalidLibraryNames[0],
			&engine.Settings.ValidStandardLibraryNames,
		)
	} else if len(invalidLibraryNames) > 1 {
		invalidLibraryConcatenation := strings.Join(invalidLibraryNames, "', '")
		return fmt.Errorf("unable to add '%s' as standard library; must be one of: %s",
			invalidLibraryConcatenation,
			&engine.Settings.ValidStandardLibraryNames,
		)
	}
	engine.Settings.StandardLibraries = libraryNames
	return nil
}

// RemoveStandardLibrary drops the specified tengo standard library from the engine's configuration, deleting it from
// the StandardLibraries setting and removing its entry from the ScriptHeader.
func (engine *Engine) RemoveStandardLibrary(libraryName string) error {
	for index, includedLibraryName := range engine.Settings.StandardLibraries {
		if includedLibraryName == libraryName {
			engine.Settings.StandardLibraries = utils.RemoveIndex(engine.Settings.StandardLibraries, index)
			scriptLines := strings.Split(engine.Settings.ScriptHeader, "\n")
			for scriptLineIndex, scriptLine := range scriptLines {
				if strings.HasPrefix(scriptLine, libraryName) {
					engine.Settings.ScriptHeader = strings.Join(utils.RemoveIndex(scriptLines, scriptLineIndex), "\n")
					break
				}
			}
			return nil
		}
	}
	return fmt.Errorf("unable to remove '%s' as standard library; not found in current list: %s",
		libraryName,
		engine.Settings.StandardLibraries,
	)
}

// AddApplicationLibrary adds the specified Tympan script library to the engine's cache and appends the declaration for
// using the library to the script header. This method otherwise works just like the AddStandardLibrary method except
// that it works on library objects and can take zero or more libraries. It is safe to call in a loop when processing a
// Tympan module for discovering and adding standalone libraries.
func (engine *Engine) AddApplicationLibraries(libraries ...Library) {
	for _, library := range libraries {
		engine.Settings.ApplicationLibraries = append(engine.Settings.ApplicationLibraries, library)
		engine.Settings.ScriptHeader += fmt.Sprintf("%s := import(\"%s\")\n", library.Name, library.Name)
	}
}

// RemoveApplicationLibrary drops the specified library from the engine's configuration, deleting it from the
// ApplicationLibraries setting and removing its entry from the ScriptHeader.
func (engine *Engine) RemoveApplicationLibrary(libraryName string) error {
	for index, includedLibrary := range engine.Settings.ApplicationLibraries {
		if includedLibrary.Name == libraryName {
			engine.Settings.ApplicationLibraries = utils.RemoveIndex(engine.Settings.ApplicationLibraries, index)
			scriptLines := strings.Split(engine.Settings.ScriptHeader, "\n")
			for scriptLineIndex, scriptLine := range scriptLines {
				if strings.HasPrefix(scriptLine, libraryName) {
					engine.Settings.ScriptHeader = strings.Join(utils.RemoveIndex(scriptLines, scriptLineIndex), "\n")
					break
				}
			}
			return nil
		}
	}
	return fmt.Errorf("unable to remove '%s' as application library; not found in current list: %s",
		libraryName,
		engine.ApplicationLibraryNames(),
	)
}

// AddApplicationModule adds the specified Tympan script module to the engine's cache and appends the declaration for
// using the module to the script header. This method otherwise works just like the AddStandardLibrary method except
// that it works on Module objects only. Note that it does not append the submodules in the script header, as the module
// should itself be the interface to any submodules. If the submodule should be available outside of its parent module,
// it should probably be included as a standalone library instead.
func (engine *Engine) AddApplicationModule(module Module) {
	engine.Settings.ApplicationModules = append(engine.Settings.ApplicationModules, module)
	engine.Settings.ScriptHeader += fmt.Sprintf("%s := import(\"%s\")\n", module.Name, module.Name)
}

// RemoveApplicationModule drops the specified module from the engine's configuration, deleting it from the
// ApplicationModules setting and removing its entry from the ScriptHeader.
func (engine *Engine) RemoveApplicationModule(moduleName string) error {
	for index, includedLibrary := range engine.Settings.ApplicationModules {
		if includedLibrary.Name == moduleName {
			engine.Settings.ApplicationModules = utils.RemoveIndex(engine.Settings.ApplicationModules, index)
			scriptLines := strings.Split(engine.Settings.ScriptHeader, "\n")
			for scriptLineIndex, scriptLine := range scriptLines {
				if strings.HasPrefix(scriptLine, moduleName) {
					engine.Settings.ScriptHeader = strings.Join(utils.RemoveIndex(scriptLines, scriptLineIndex), "\n")
					break
				}
			}
			return nil
		}
	}
	return fmt.Errorf("unable to remove '%s' as application module; not found in current list: %s",
		moduleName,
		engine.ApplicationModuleNames(),
	)
}

// Adds a new script to the engine from a given name and script body as string. At the time the script is added, all
// necessary actions are taken to ensure the script can be run immediately after. This means that you want to be sure
// to configure the engine with desired libraries and settings before adding any scripts. The script is not
// automatically compiled on creation, so it is possible to update/set new variables and reuse the script.
func (engine *Engine) AddScript(name string, scriptString string) error {
	if engine.ScriptWithSameNameExists(name) {
		return fmt.Errorf("cannot add script '%s' to engine: script with the same name already exists", name)
	}
	scriptWithHeader := strings.Join([]string{engine.Settings.ScriptHeader, scriptString}, "\n\n")
	script := tengo.NewScript([]byte(scriptWithHeader))
	engine.InitializeLibraryImporter()
	script.SetImports(engine.Importer)
	engine.Scripts = append(engine.Scripts, &Script{Name: name, Body: scriptWithHeader, Script: script})
	return nil
}

// Retrieves a cached script from the engine by name.
func (engine *Engine) GetScript(name string) *Script {
	for _, script := range engine.Scripts {
		if script.Name == name {
			return script
		}
	}
	return nil
}

// Checks to see whether the specified name matches a script the Engine has already cached.
func (engine *Engine) ScriptWithSameNameExists(name string) bool {
	for _, script := range engine.Scripts {
		if script.Name == name {
			return true
		}
	}
	return false
}
