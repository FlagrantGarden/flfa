package flfa_test

import (
	"fmt"
	"testing"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/spf13/afero"
)

func TestApi_CacheModuleData(t *testing.T) {
	tests := []struct {
		name             string
		moduleFolderPath string
	}{
		{
			name:             "shit",
			moduleFolderPath: "../../modules/core",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewOsFs() // use the real file system
			api := &flfa.Api{
				Tympan: &tympan.Tympan{
					AFS:  &afero.Afero{Fs: fs},
					IOFS: &afero.IOFS{Fs: fs},
				},
			}
			api.CacheModuleData(tt.moduleFolderPath)
			fmt.Printf("oh snaaaaap")
		})
	}
}
