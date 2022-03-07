package flfa

import (
	"os"
	"path/filepath"
)

type State struct {
	SavePath       string
	Guid           string
	FileInfo       os.FileInfo
	ActiveSkirmish string
	Skirmishes     []string
}

func (ffapi *Api) InitializeState(path string) (State, error) {
	return State{}, nil
}

func (s *State) Name() string {
	return filepath.Base(s.SavePath)
}
