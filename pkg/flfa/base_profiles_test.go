package flfa

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
)

func TestApi_ReadAndParseProfileData(t *testing.T) {
	tests := []struct {
		name         string
		dataFilePath string
		want         []Trait
	}{
		{
			name:         "fuck",
			dataFilePath: "../../modules/core/Profiles.yaml",
			want:         []Trait{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewOsFs() // use the real file system
			ffapi := &Api{
				AFS:          &afero.Afero{Fs: fs},
				IOFS:         &afero.IOFS{Fs: fs},
				CachedTraits: []Trait{},
			}
			results, err := ffapi.ReadAndParseProfileData(tt.dataFilePath)
			if err != nil {
				t.Errorf("fuck!")
			} else {
				fmt.Printf("Results! %+v", results)
			}
		})
	}
}
