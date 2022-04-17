package flfa

import (
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/rs/zerolog/log"
)

func (ffapi *Api) InitializeGameState() error {
	err := ffapi.Tympan.InitializeConfig()
	if err != nil {
		return err
	}

	log.Trace().Msgf("Loading module data from %s", ffapi.Tympan.Configuration.FolderPaths.Cache)
	installedModules, err := ffapi.InstalledModules()
	if err != nil {
		log.Error().Msgf("error initializing game; unable to list installed modules: %s", err)
	}
	log.Trace().Msgf("Installed modules: %s", strings.Join(installedModules, ", "))
	if !utils.Contains(installedModules, "core") {
		ffapi.CacheModuleData("modules/core", true)
	}
	for _, module := range installedModules {
		ffapi.CacheModuleData(filepath.Join(ffapi.Tympan.Configuration.FolderPaths.Application, "modules", module), false)
	}
	log.Trace().Msgf("Caching personas from %s", ffapi.Tympan.Configuration.FolderPaths.Cache)
	ffapi.CachePlayers("")

	ffapi.InitializeEngine()
	return nil
}
