package editor

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/editor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type EditorCommand struct {
	Api *flfa.Api
}

type EditorCommander interface {
	CreateCommand() *cobra.Command
}

func (e *EditorCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "editor",
		Short:             "Manage custom content",
		Long:              "Manage custom content by creating, editing, and deleting companies, groups, and rosters",
		PersistentPreRunE: e.initialize,
		RunE:              e.execute,
	}

	cmd.Flags().SortFlags = false

	// if human readable, return terminal display
	// if json, write json blob
	// if export, write markdown
	// if persona specified, skip persona chooser
	// if company specified, skip company chooser

	return cmd
}

func (e *EditorCommand) initialize(cmd *cobra.Command, args []string) error {
	return e.Api.InitializeGameState()
}

func (e *EditorCommand) execute(cmd *cobra.Command, args []string) error {
	model := editor.NewModel(e.Api)
	program := tea.NewProgram(model, tea.WithAltScreen())
	return program.Start()
}
