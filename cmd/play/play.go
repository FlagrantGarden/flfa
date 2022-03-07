package play

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/terminal_documentation"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PlayCommand struct {
	Api          *flfa.Api
	SaveFilePath string
}

type PlayCommander interface {
	CreateCommand() *cobra.Command
}

func (p *PlayCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "play",
		Short:             "Play the game",
		Long:              "Play the game",
		PersistentPreRunE: p.initializeGameState,
		// Args: handleArgs,
		// ValidArgsFunction: flagCompletion,
		// PreRun: preExecute,
		RunE: p.execute,
	}
	cmd.Flags().SortFlags = false
	// cmd.Flags.BoolVarP(&p.List)
	return cmd
}

func (p *PlayCommand) initializeGameState(cmd *cobra.Command, args []string) error {
	if p.SaveFilePath == "" {
		p.SaveFilePath = getSaveFilePath()
	}
	log.Trace().Msgf("Using save file at: %s", p.SaveFilePath)
	log.Trace().Msgf("Loading module data from %s", p.Api.Tympan.RunningConfig.ModulePath)
	p.Api.CacheModuleData(filepath.Join(p.Api.Tympan.RunningConfig.ModulePath, "core"))
	return nil
}

func getSaveFilePath() string {
	filename := fmt.Sprintf("%s.yaml", viper.GetString(tympan.CurrentGameKey))
	return filepath.Join(viper.GetString(tympan.UserDataPathKey), filename)
}

func (p *PlayCommand) execute(cmd *cobra.Command, args []string) error {
	group, err := p.Api.NewGroupPrompt()
	if err != nil {
		return err
	}

	switch viper.GetString("format") {
	case "json":
		fmt.Print(group.JSON())
	default:
		groupOutput := strings.Builder{}
		groupOutput.WriteString("| Name | Profile | Melee | Missile | Move | FS | R | T | Traits |\n")
		groupOutput.WriteString("| ---- | ------- | ----- | ------- | ---- | -- | - | - | ------ |\n")
		groupOutput.WriteString(group.MarkdownTableEntry())
		td := terminal_documentation.TerminalDocumentation{}
		output, err := td.Render(groupOutput.String())
		if err != nil {
			return err
		}
		fmt.Print(output)
	}
	return nil
}
