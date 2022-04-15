package scripting

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// GetModule returns a Module instance with all of its submodules. It requires a path to the root folder of a Tympan
// module and an afero file system. It looks in the absolute path to the module folder for a "scripts" folder and a
// ".tengo" file that is named the same as the module folder. If it finds the file, it uses GetLibrary to retrieve it,
// setting the Module's Name and Body to the appropriate values. If it does not find the file, it returns immediately.
//
// If the module file is found and retrieved without error, it then looks for the "submodules" folder in the same
// directory as the module script file. If that directory exists, it calls GetFolderLibraries on it and adds all of the
// discovered libraries to the Module's Submodules list.
func GetModule(moduleFolderPath string, afs *afero.Afero) (module Module, err error) {
	absModuleFolderPath, err := filepath.Abs(moduleFolderPath)
	if err != nil {
		return module, fmt.Errorf("unable to determine absolute path to module folder '%s': %s", moduleFolderPath, err)
	}

	moduleName := filepath.Base(moduleFolderPath)
	moduleFilePath := filepath.Join(absModuleFolderPath, "scripts", fmt.Sprintf("%s.tengo", moduleName))

	// find the module
	exists, err := afs.Exists(moduleFilePath)
	if err != nil {
		return module, fmt.Errorf("unable to determine if script module file '%s' exists: %s", moduleFilePath, err)
	} else if exists {
		module.Library, err = GetLibrary(moduleFilePath, afs)
		if err != nil {
			return
		}
	} else {
		return
	}

	// add the submodules to module
	submoduleFolderPath := filepath.Join(absModuleFolderPath, "scripts", "submodules")
	exists, err = afs.DirExists(submoduleFolderPath)
	if err != nil {
		return module, fmt.Errorf("unable to determine if submodule folder '%s' exists: %s", moduleFolderPath, err)
	} else if exists {
		module.Submodules, err = GetFolderLibraries(submoduleFolderPath, afs)
	}

	return
}

// GetStandaloneLibraries is a helper function for returning the list of all libraries found in a Tympan module folder.
// It requires the path to the root folder of a Tympan module and an Afero file system. It looks in the absolute path to
// the module folder for the "scripts" folder with a "libraries" subfolder. If that folder exists, it calls
// GetFolderLibraries on it to return all of the standalone libraries the module provides.
func GetStandaloneLibraries(moduleFolderPath string, afs *afero.Afero) (libraries []Library, err error) {
	absModuleFolderPath, err := filepath.Abs(moduleFolderPath)
	if err != nil {
		return libraries, fmt.Errorf("unable to determine absolute path to module folder '%s': %s", moduleFolderPath, err)
	}

	libraryFolderPath := fmt.Sprintf("%s/scripts/libraries", absModuleFolderPath)
	exists, err := afs.DirExists(libraryFolderPath)
	if err != nil {
		return libraries, fmt.Errorf("unable to determine if standalone library folder '%s' exists: %s", moduleFolderPath, err)
	} else if exists {
		return GetFolderLibraries(libraryFolderPath, afs)
	}

	return
}

// GetFolderLibraries requires the path to a folder containing *.tengo files you want to add as libraries and an afero
// file system to use. It looks for the absolute path to the folder and then walks it, calling GetLibrary on each tengo
// file it finds, appending found libraries to the list of libraries to return in the order they're found.
//
// If any errors occur while walking the folder, it stops looking for more libraries and returns the successfully parsed
// libraries and the error that stopped the execution.
func GetFolderLibraries(folderPath string, afs *afero.Afero) (libraries []Library, err error) {
	folderPath, err = filepath.Abs(folderPath)
	if err != nil {
		return libraries, fmt.Errorf("unable to determine absolute path to script library file '%s': %s", folderPath, err)
	}
	// search the libraries folder, add each found library
	err = afs.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		isScriptLibrary, _ := filepath.Match("*.tengo", filepath.Base(path))
		if isScriptLibrary {
			library, err := GetLibrary(path, afs)
			if err != nil {
				return err
			}
			libraries = append(libraries, library)
		}
		return nil
	})
	return
}

// GetLibrary requires the path to a *.tengo file you want to add as a library and an afero file system to use. It looks
// for the absolute path to the file, tries to read it, and returns a Library (with the Name set to the file's name
// -- without the ".tengo" extension -- and the Body set to the contents of the file) and nil for the error.
//
// If the file can't be read for any reason, it returns an empty Library and the error.
func GetLibrary(filePath string, afs *afero.Afero) (library Library, err error) {
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return library, fmt.Errorf("unable to determine absolute path to script library file '%s': %s", filePath, err)
	}

	contents, err := afs.ReadFile(filePath)
	if err != nil {
		return library, fmt.Errorf("unable to read script library file '%s': %s", filePath, err)
	}

	library.Name = strings.TrimSuffix(filepath.Base(filePath), ".tengo")
	library.Body = string(contents)

	return
}

// GetModule returns a Module instance with all of its submodules. It requires a path to the root folder of a Tympan
// module and an embeddedd file system. It looks in the module folder for a "scripts" folder and a ".tengo" file that is
// named the same as the module folder. If it finds the file, it uses GetLibrary to retrieve it, setting the Module's
// Name and Body to the appropriate values. If it does not find the file, it returns immediately.
//
// If the module file is found and retrieved without error, it then looks for the "submodules" folder in the same
// directory as the module script file, calling GetEmbeddedFolderLibraries on it and adding all of the discovered
// libraries to the Module's Submodules list.
func GetEmbeddedModule(folderPath string, efs *embed.FS) (module Module, err error) {
	moduleName := filepath.Base(folderPath)
	moduleFilePath := fmt.Sprintf("%s/scripts/%s.tengo", folderPath, moduleName)

	// find the module
	module.Library, err = GetEmbeddedLibrary(moduleFilePath, efs)
	if err != nil || module.Name == "" {
		return
	}

	// add the submodules to module
	submoduleFolderPath := fmt.Sprintf("%s/scripts/submodules", folderPath)
	module.Submodules, err = GetEmbeddedFolderLibraries(submoduleFolderPath, efs)

	return
}

// GetEmbeddedStandaloneLibraries is a helper function for returning the list of all libraries found in a Tympan module
// folder. It requires the path to the root folder of a Tympan module and an embedded file system. It figures out the
// path to the "libraries" subfolder inside the "scripts" folder of the specified module path and then calls
// GetEmbeddedFolderLibraries on it to return all of the standalone libraries the module provides.
func GetEmbeddedStandaloneLibraries(folderPath string, efs *embed.FS) (libraries []Library, err error) {
	libraryFolderPath := fmt.Sprintf("%s/scripts/libraries", folderPath)

	return GetEmbeddedFolderLibraries(libraryFolderPath, efs)
}

// GetEmbeddedFolderLibraries requires the path to a folder containing *.tengo files you want to add as libraries and an
// embedded file system to use. It walks the specified folder, calling GetLibrary on each tengo file it finds, appending
// found libraries to the list of libraries to return in the order they're found.
//
// If any errors occur while walking the folder, it stops looking for more libraries and returns the successfully parsed
// libraries and the error that stopped the execution.
func GetEmbeddedFolderLibraries(folderPath string, efs *embed.FS) (libraries []Library, err error) {
	// search the libraries folder, add each found library
	err = fs.WalkDir(efs, folderPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		isScriptLibrary, _ := filepath.Match("*.tengo", filepath.Base(path))
		if isScriptLibrary {
			library, err := GetEmbeddedLibrary(path, efs)
			if err != nil {
				return err
			}
			libraries = append(libraries, library)
		}
		return nil
	})
	return
}

// GetEmbeddedLibrary requires the path to a *.tengo file you want to add as a library and an embedded file system to
// use. It looks for the specified path to the file, tries to read it, and returns a Library (with the Name set to the
// file's name -- without the ".tengo" extension -- and the Body set to the contents of the file) and nil for the error.
//
// If the file can't be read for any reason, it returns an empty Library and the error.
func GetEmbeddedLibrary(filePath string, efs *embed.FS) (library Library, err error) {
	contents, err := efs.ReadFile(filePath)
	if err != nil {
		return library, fmt.Errorf("unable to read script library file '%s': %s", filePath, err)
	}

	library.Name = strings.TrimSuffix(filepath.Base(filePath), ".tengo")
	library.Body = string(contents)

	return
}
