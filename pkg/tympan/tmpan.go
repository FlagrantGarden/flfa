package tympan

import (
	"github.com/Masterminds/semver"
	"github.com/spf13/afero"
)

type Tympan struct {
	AFS           *afero.Afero
	IOFS          *afero.IOFS
	RunningConfig Config
}

type Tympaner interface {
	GenerateDefaultConfig()
	LoadConfig() (err error)
	GetDefaultModulePath() (string, error)
	GetDefaultUserDataPath() (string, error)
	GetDefaultApiVersion() *semver.Version
}
