package flfa

import (
	"embed"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/scripting"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/skirmish"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/user"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/FlagrantGarden/flfa/pkg/tympan/module"
	tympan_scripting "github.com/FlagrantGarden/flfa/pkg/tympan/module/scripting"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/instance"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
)

type Api struct {
	Tympan       *tympan.Tympan[*Configuration]
	EMFS         *embed.FS
	Cache        DataCache
	ScriptEngine *tympan_scripting.Engine
}

type DataCache struct {
	Traits          []data.Trait
	Profiles        []data.Profile
	Spells          []data.Spell
	Companies       []data.Company
	Personas        []*persona.Persona[user.Data, user.Settings]
	ScriptModules   []tympan_scripting.Module
	ScriptLibraries []tympan_scripting.Library
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
		ffapi.ScriptEngine = scripting.NewEngine(ffapi.Cache.ScriptModules, ffapi.Cache.ScriptLibraries)
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
	userKind := &state.Kind{
		Name:       "user",
		FolderName: "users",
	}
	discoveredPersonas, _ := persona.DiscoverPersonas[user.Data, user.Settings](userKind, cachePath, ffapi.Tympan.AFS)
	ffapi.Cache.Personas = discoveredPersonas
}

func (ffapi *Api) GetUserPersona(name string, cachePath string) (*persona.Persona[user.Data, user.Settings], error) {
	if cachePath == "" {
		cachePath = ffapi.Tympan.Configuration.FolderPaths.Cache
	}
	userKind := user.Kind()
	return persona.GetPersona[user.Data, user.Settings](name, userKind, cachePath, ffapi.Tympan.AFS)
}

func (ffapi *Api) GetActiveSkirmish(activeUserPersona *persona.Persona[user.Data, user.Settings], cachePath string) (*instance.Instance[skirmish.Skirmish], error) {
	if cachePath == "" {
		cachePath = ffapi.Tympan.Configuration.FolderPaths.Cache
	}
	skirmishKind := skirmish.Kind()
	skirmishPersona := &instance.Persona{
		Name: activeUserPersona.Name,
		Kind: activeUserPersona.Kind,
	}
	return instance.GetInstance[skirmish.Skirmish](activeUserPersona.Settings.ActiveSkirmish, skirmishKind, skirmishPersona, cachePath, ffapi.Tympan.AFS)
}
