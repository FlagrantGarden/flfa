package flfa

import (
	"embed"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/player"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/skirmish"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/FlagrantGarden/flfa/pkg/tympan/module/scripting"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/instance"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
)

type Api struct {
	Tympan       *tympan.Tympan[*Configuration]
	EMFS         *embed.FS
	Cache        DataCache
	ScriptEngine *scripting.Engine
}

type DataCache struct {
	Traits          []data.Trait
	Profiles        []data.Profile
	Spells          []data.Spell
	Companies       []data.Company
	Personas        []player.Player
	ScriptModules   []scripting.Module
	ScriptLibraries []scripting.Library
}

type Configuration struct {
	tympan.SharedConfig `mapstructure:",squash" tympanconfig:"ignore"`
	ActiveUserPersona   string `mapstructure:"active_user_persona"`
}

func (config *Configuration) Initialize() error {
	return nil
}

func (ffapi *Api) InitializeEngine() {
	if ffapi.ScriptEngine == nil {
		ffapi.ScriptEngine = scripting.NewEngine()
		// ignore errors for now
		ffapi.ScriptEngine.SetStandardLibraries(ffapi.ScriptEngine.AllowedStandardLibraries())
		ffapi.ScriptEngine.AddApplicationLibraries(ffapi.Cache.ScriptLibraries...)
		for _, module := range ffapi.Cache.ScriptModules {
			ffapi.ScriptEngine.AddApplicationModule(module)
		}
	}
}

func (ffapi *Api) CacheModuleData(modulePath string, embedded bool) {
	ffapi.CacheProfiles(modulePath, embedded)
	ffapi.CacheTraits(modulePath, embedded)
	ffapi.CacheSpells(modulePath, embedded)
	ffapi.CacheCompanies(modulePath, embedded)
	ffapi.CacheScriptLibraries(modulePath, embedded)
	ffapi.CacheScriptModules(modulePath, embedded)
}

func (ffapi *Api) InstalledModules() (installedModules []string, err error) {
	moduleFolderPath := filepath.Join(ffapi.Tympan.Configuration.FolderPaths.Cache, "modules")

	moduleFolderExists, err := ffapi.Tympan.AFS.DirExists(moduleFolderPath)
	if err != nil || !moduleFolderExists {
		return
	}

	cacheFolderItems, err := ffapi.Tympan.AFS.ReadDir(moduleFolderPath)
	if err != nil {
		return
	}

	for _, item := range cacheFolderItems {
		if item.IsDir() {
			installedModules = append(installedModules, item.Name())
		}
	}

	return
}

func (ffapi *Api) CacheUserPersonas(cachePath string) {
	if cachePath == "" {
		cachePath = ffapi.Tympan.Configuration.FolderPaths.Cache
	}
	discoveredPersonas, _ := persona.DiscoverPersonas[player.Data, player.Settings](player.Kind(), cachePath, ffapi.Tympan.AFS)
	for _, persona := range discoveredPersonas {
		ffapi.Cache.Personas = append(ffapi.Cache.Personas, player.Player{Persona: persona})
	}
}

func (ffapi *Api) GetUserPersona(name string, cachePath string) (*persona.Persona[player.Data, player.Settings], error) {
	if cachePath == "" {
		cachePath = ffapi.Tympan.Configuration.FolderPaths.Cache
	}
	return persona.GetPersona[player.Data, player.Settings](name, player.Kind(), cachePath, ffapi.Tympan.AFS)
}

func (ffapi *Api) GetActiveSkirmish(activeUserPersona *persona.Persona[player.Data, player.Settings], cachePath string) (*instance.Instance[skirmish.Skirmish], error) {
	if cachePath == "" {
		cachePath = ffapi.Tympan.Configuration.FolderPaths.Cache
	}
	skirmishPersona := &instance.Persona{
		Name: activeUserPersona.Name,
		Kind: activeUserPersona.Kind,
	}
	return instance.GetInstance[skirmish.Skirmish](activeUserPersona.Settings.ActiveSkirmish, skirmish.Kind(), skirmishPersona, cachePath, ffapi.Tympan.AFS)
}
