package flfa

import (
	"embed"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/skirmish"
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/user"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/FlagrantGarden/flfa/pkg/tympan/module"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/instance"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
)

type Api struct {
	Tympan *tympan.Tympan[*Configuration]
	EMFS   *embed.FS
	Cache  DataCache
}

type DataCache struct {
	Traits    []data.Trait
	Profiles  []data.Profile
	Spells    []data.Spell
	Companies []data.Company
	Personas  []*persona.Persona[user.Data, user.Settings]
}

type Configuration struct {
	tympan.SharedConfig `mapstructure:",squash" tympanconfig:"ignore"`
	ActiveUserPersona   string `mapstructure:"active_user_persona"`
}

func (config *Configuration) Initialize() error {
	return nil
}

func (ffapi *Api) CacheModuleData(modulePath string, embedded bool) {
	// load profiles
	var profiles []data.Profile
	if embedded {
		profiles, _ = module.GetEmbeddedDataByFile[data.Profile](modulePath, "Profiles", ffapi.EMFS)
	} else {
		profiles, _ = module.GetDataByFile[data.Profile](modulePath, "Profiles", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Profiles = append(ffapi.Cache.Profiles, profiles...)
	// load traits
	var traits []data.Trait
	if embedded {
		traits, _ = module.GetEmbeddedDataByFolder[data.Trait](modulePath, "Traits", ffapi.EMFS)
	} else {
		traits, _ = module.GetDataByFolder[data.Trait](modulePath, "Traits", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Traits = append(ffapi.Cache.Traits, traits...)
	// load spells
	var spells []data.Spell
	if embedded {
		spells, _ = module.GetEmbeddedDataByFile[data.Spell](modulePath, "Spells", ffapi.EMFS)
	} else {
		spells, _ = module.GetDataByFile[data.Spell](modulePath, "Spells", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Spells = append(ffapi.Cache.Spells, spells...)
	// load companies
	var companies []data.Company
	if embedded {
		companies, _ = module.GetEmbeddedDataByFile[data.Company](modulePath, "Companies", ffapi.EMFS)
	} else {
		companies, _ = module.GetDataByFile[data.Company](modulePath, "Companies", ffapi.Tympan.AFS)
	}
	for _, company := range companies {
		company.Initialize(profiles, traits)
		ffapi.Cache.Companies = append(ffapi.Cache.Companies, company)
	}
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
