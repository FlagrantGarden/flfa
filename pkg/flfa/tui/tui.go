package tui

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/tympan/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

type SharedModel struct {
	Api *flfa.Api
	compositor.Compositor
}

func (model *SharedModel) SaveConfig() tea.Cmd {
	return func() tea.Msg {
		err := model.Api.Tympan.SaveConfig()
		if err != nil {
			return model.RecordFatalError(err)
		}
		model.State = compositor.StateSavedConfiguration
		return nil
	}
}
