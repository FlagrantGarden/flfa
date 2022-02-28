package play

import (
	"fmt"
	"path/filepath"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
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
	log.Trace().Msgf("Loading module data from %s", p.Api.RunningConfig.ModulePath)
	p.Api.CacheModuleData(filepath.Join(p.Api.RunningConfig.ModulePath, "core"))
	return nil
}

func getSaveFilePath() string {
	filename := fmt.Sprintf("%s.yaml", viper.GetString(flfa.CurrentGameKey))
	return filepath.Join(viper.GetString(flfa.UserDataPathKey), filename)
}

func (p *PlayCommand) execute(cmd *cobra.Command, args []string) error {
	group, err := p.Api.NewGroupPrompt()
	if err != nil {
		return err
	}
	fmt.Printf("Your group is: %+v", group)
	return nil
}
