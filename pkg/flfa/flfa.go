package flfa

import (
	"github.com/spf13/afero"
)

type Api struct {
	AFS             *afero.Afero
	IOFS            *afero.IOFS
	RunningConfig   Config
	CachedTraits    []Trait
	CachedProfiles  []BaseProfile
	CachedSpells    []Spell
	CachedCompanies []Company
}

func (ffapi *Api) CacheModuleData(path string) {
	// load base profiles
	ffapi.CacheBaseProfiles(path)
	// load traits
	ffapi.CacheModuleTraits(path)
	// load spells
	ffapi.CacheModuleSpells(path)
	// load companies
	ffapi.CacheModuleCompanies(path)
}
