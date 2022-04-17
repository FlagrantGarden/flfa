package player

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/data"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/persona"
)

type Player struct {
	*persona.Persona[Data, Settings]
}

type Settings struct {
	ActiveSkirmish string `mapstructure:"active_skirmish"`
	Skirmishes     []Skirmish
}

type Skirmish struct {
	Name          string
	Configuration SkirmishConfiguration
}

type SkirmishConfiguration struct {
	Autosave bool
}

func (playerSettings Settings) Initialize() *Settings {
	// Check if empty; for now, the implementation is such that UserSettings should always have an ActiveSkirmish,
	// so just verify that it isn't nil and, if it is, create the struct and initialize it.
	if playerSettings.ActiveSkirmish != "" {
		return &playerSettings
	}
	return &Settings{
		ActiveSkirmish: "default",
		Skirmishes: []Skirmish{
			{
				Name: "default",
				Configuration: SkirmishConfiguration{
					Autosave: true,
				},
			},
		},
	}
}

type Data struct {
	Companies []data.Company
}

func (playerData Data) Initialize() *Data {
	// check if UserData is empty; for now it only has one field, so just check length
	if len(playerData.Companies) != 0 {
		return &playerData
	}
	return &Data{}
}

func Kind() *state.Kind {
	return &state.Kind{
		Name:       "player",
		FolderName: "players",
	}
}
