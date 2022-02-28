package flfa

import (
	"testing"

	"github.com/spf13/afero"
)

func TestApi_CacheModuleCompanies(t *testing.T) {
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
			fs := afero.NewOsFs() // use the real file system
			ffapi := &Api{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
			}
			ffapi.CacheBaseProfiles(tt.modulePath)
			ffapi.CacheModuleTraits(tt.modulePath)

			if err := ffapi.CacheModuleCompanies(tt.modulePath); (err != nil) != tt.wantErr {
				t.Errorf("Api.CacheModuleCompanies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
