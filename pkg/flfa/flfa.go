package flfa

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
)

type Api struct {
	Tympan          *tympan.Tympan
	CachedTraits    []data.Trait
	CachedProfiles  []data.Profile
	CachedSpells    []data.Spell
	CachedCompanies []data.Company
}

func (ffapi *Api) CacheModuleData(modulePath string) {
	// load profiles
	profiles, _ := tympan.GetModuleDataByFile[data.Profile](modulePath, "Profiles", ffapi.Tympan.AFS)
	ffapi.CachedProfiles = append(ffapi.CachedProfiles, profiles...)
	// load traits
	traits, _ := tympan.GetModuleDataByFolder[data.Trait](modulePath, "Traits", ffapi.Tympan.AFS)
	ffapi.CachedTraits = append(ffapi.CachedTraits, traits...)
	// load spells
	spells, _ := tympan.GetModuleDataByFile[data.Spell](modulePath, "Spells", ffapi.Tympan.AFS)
	ffapi.CachedSpells = append(ffapi.CachedSpells, spells...)
	// load companies
	companies, _ := tympan.GetModuleDataByFile[data.Company](modulePath, "Companies", ffapi.Tympan.AFS)
	for _, company := range companies {
		company.Initialize(profiles, traits)
		ffapi.CachedCompanies = append(ffapi.CachedCompanies, company)
	}
}
