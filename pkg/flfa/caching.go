package flfa

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/module"
	tympan_scripting "github.com/FlagrantGarden/flfa/pkg/tympan/module/scripting"
	"github.com/spf13/afero"
)

func (ffapi *Api) CachingFs(embedded bool) *afero.Afero {
	afs := ffapi.Tympan.AFS
	if embedded {
		afs = &afero.Afero{
			Fs: afero.FromIOFS{FS: ffapi.EMFS},
		}
	}
	return afs
}

func (ffapi *Api) CacheProfiles(modulePath string, embedded bool) {
	var profiles []data.Profile
	if embedded {
		profiles, _ = module.GetEmbeddedDataByFile[data.Profile](modulePath, "Profiles", ffapi.EMFS)
	} else {
		profiles, _ = module.GetDataByFile[data.Profile](modulePath, "Profiles", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Profiles = append(ffapi.Cache.Profiles, profiles...)
}

func (ffapi *Api) CacheTraits(modulePath string, embedded bool) {
	var traits []data.Trait
	if embedded {
		traits, _ = module.GetEmbeddedDataByFolder[data.Trait](modulePath, "Traits", ffapi.EMFS)
	} else {
		traits, _ = module.GetDataByFolder[data.Trait](modulePath, "Traits", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Traits = append(ffapi.Cache.Traits, traits...)
}

func (ffapi *Api) CacheSpells(modulePath string, embedded bool) {
	var spells []data.Spell
	if embedded {
		spells, _ = module.GetEmbeddedDataByFile[data.Spell](modulePath, "Spells", ffapi.EMFS)
	} else {
		spells, _ = module.GetDataByFile[data.Spell](modulePath, "Spells", ffapi.Tympan.AFS)
	}
	ffapi.Cache.Spells = append(ffapi.Cache.Spells, spells...)
}

func (ffapi *Api) CacheCompanies(modulePath string, embedded bool) {
	var companies []data.Company
	if embedded {
		companies, _ = module.GetEmbeddedDataByFile[data.Company](modulePath, "Companies", ffapi.EMFS)
	} else {
		companies, _ = module.GetDataByFile[data.Company](modulePath, "Companies", ffapi.Tympan.AFS)
	}
	for _, company := range companies {
		company.Initialize(ffapi.Cache.Profiles, ffapi.Cache.Traits)
		ffapi.Cache.Companies = append(ffapi.Cache.Companies, company)
	}
}

func (ffapi *Api) CacheRosters(modulePath string, embedded bool) {}

func (ffapi *Api) CacheScriptLibraries(modulePath string, embedded bool) {
	var scriptLibraries []tympan_scripting.Library
	if embedded {
		scriptLibraries, _ = tympan_scripting.GetEmbeddedStandaloneLibraries(modulePath, ffapi.EMFS)
	} else {
		scriptLibraries, _ = tympan_scripting.GetStandaloneLibraries(modulePath, ffapi.Tympan.AFS)
	}
	ffapi.Cache.ScriptLibraries = append(ffapi.Cache.ScriptLibraries, scriptLibraries...)
}

func (ffapi *Api) CacheScriptModules(modulePath string, embedded bool) {
	var scriptModule tympan_scripting.Module
	if embedded {
		scriptModule, _ = tympan_scripting.GetEmbeddedModule(modulePath, ffapi.EMFS)
	} else {
		scriptModule, _ = tympan_scripting.GetModule(modulePath, ffapi.Tympan.AFS)
	}
	ffapi.Cache.ScriptModules = append(ffapi.Cache.ScriptModules, scriptModule)
}
