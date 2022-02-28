package flfa

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs() // use the real file system
var ffapi = &Api{
	AFS:          &afero.Afero{Fs: fs},
	IOFS:         &afero.IOFS{Fs: fs},
	CachedTraits: []Trait{},
}

func TestFlfa_ReadAndParseTraitData(t *testing.T) {
	tests := []struct {
		name         string
		dataFilePath string
		want         []Trait
	}{
		{
			name:         "fuck",
			dataFilePath: "../../modules/core/traits/captain.yaml",
			want:         []Trait{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ffapi.ReadAndParseTraitData(tt.dataFilePath)
			if err != nil {
				t.Errorf("fuck!")
			} else {
				fmt.Printf("Results! %+v", results)
			}
		})
	}
}

func TestApi_CacheModuleTraits(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		wantErr    bool
	}{
		{
			name:       "shit",
			modulePath: "../../modules/core",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ffapi.CacheModuleTraits(tt.modulePath); (err != nil) != tt.wantErr {
				t.Errorf("Api.CacheModuleTraits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
