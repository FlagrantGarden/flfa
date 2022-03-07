package tympan_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/spf13/afero"
)

func TestReadAndParseData(t *testing.T) {
	tests := []struct {
		name         string
		dataFilePath string
	}{
		{
			name:         "fuck",
			dataFilePath: "../../modules/core/Profiles.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewOsFs() // use the real file system
			tapi := &tympan.Tympan{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
			}
			results, err := tympan.ReadAndParseData[flfa.Profile](tt.dataFilePath, tapi.AFS)
			if err != nil {
				t.Errorf("fuck!")
			} else {
				fmt.Printf("Results! %+v", results)
			}
		})
	}
}

func TestGetModuleDataByFile(t *testing.T) {
	tests := []struct {
		name             string
		moduleFolderPath string
		dataTypeName     string
		dataType         reflect.Type
	}{
		{
			name:             "retrieving profiles",
			moduleFolderPath: "../../modules/core",
			dataTypeName:     "Profiles",
		},
		{
			name:             "retrieving spells",
			moduleFolderPath: "../../modules/core",
			dataTypeName:     "Spells",
		},
		{
			name:             "retrieving companies",
			moduleFolderPath: "../../modules/core",
			dataTypeName:     "Companies",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewOsFs() // use the real file system
			api := &flfa.Api{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
			}
			switch tt.dataTypeName {
			case "Profiles":
				results, err := tympan.GetModuleDataByFile[flfa.Profile](tt.moduleFolderPath, tt.dataTypeName, api.AFS)
				if err != nil {
					t.Errorf("fuck!")
				} else {
					fmt.Printf("Results! %+v", results)
				}
			case "Spells":
				results, err := tympan.GetModuleDataByFile[flfa.Spell](tt.moduleFolderPath, tt.dataTypeName, api.AFS)
				if err != nil {
					t.Errorf("fuck!")
				} else {
					fmt.Printf("Results! %+v", results)
				}
			case "Companies":
				results, err := tympan.GetModuleDataByFile[flfa.Company](tt.moduleFolderPath, tt.dataTypeName, api.AFS)
				if err != nil {
					t.Errorf("fuck!")
				} else {
					fmt.Printf("Results! %+v", results)
				}
			}
		})
	}
}

func TestGetModuleDataByFolder(t *testing.T) {
	tests := []struct {
		name             string
		moduleFolderPath string
		dataFolderName   string
	}{
		{
			name:             "shit",
			moduleFolderPath: "../../modules/core",
			dataFolderName:   "Traits",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewOsFs() // use the real file system
			api := &flfa.Api{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
			}
			switch tt.dataFolderName {
			case "Traits":
				results, err := tympan.GetModuleDataByFolder[flfa.Trait](tt.moduleFolderPath, tt.dataFolderName, api.AFS)
				if err != nil {
					t.Errorf("fuck!")
				} else {
					fmt.Printf("Results! %+v", results)
				}
			}
		})
	}
}
