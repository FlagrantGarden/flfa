package play

import (
	"fmt"
	"os"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/company"
	"github.com/FlagrantGarden/flfa/pkg/flfa/tui/persona"
	"github.com/FlagrantGarden/flfa/pkg/tympan/dossier"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type PlayCommand struct {
	Api     *flfa.Api
	Dossier *dossier.Dossier
}

type PlayCommander interface {
	CreateCommand() *cobra.Command
}

func (p *PlayCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "play",
		Short:             "Play the game",
		Long:              "Play the game",
		PersistentPreRunE: p.initialize,
		// Args: handleArgs,
		// ValidArgsFunction: flagCompletion,
		// PreRun: preExecute,
		RunE: p.execute,
	}
	cmd.Flags().SortFlags = false
	// cmd.Flags.BoolVarP(&p.List)
	return cmd
}

func (p *PlayCommand) initialize(cmd *cobra.Command, args []string) error {
	return p.Api.InitializeGameState()
}

func (p *PlayCommand) execute(cmd *cobra.Command, args []string) error {
	personaModel := persona.NewModel(p.Api)
	personaProgram := tea.NewProgram(personaModel)
	if err := personaProgram.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	// c := p.Api.Cache.Companies[0]
	// fmt.Print(data.DisplayGroupTerminal(c.Groups...))
	// groupModel := group.NewModel(p.Api)
	// groupProgram := tea.NewProgram(groupModel, tea.WithAltScreen())
	// if err := groupProgram.Start(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v", err)
	// 	os.Exit(1)
	// }

	companyModel := company.NewModel(p.Api)
	companyProgram := tea.NewProgram(companyModel, tea.WithAltScreen())
	if err := companyProgram.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	// if p.Api.Tympan.Configuration.ActiveUserPersona == "" {
	// 	err := prompts.FirstPlay(p.Api)
	// 	if err != nil {
	// 		log.Logger.Error().Msgf("Something broke: %s", err)
	// 	}
	// } else {
	// 	// keep going!
	// 	activeUserPersona, err := p.Api.GetUserPersona(p.Api.Tympan.Configuration.ActiveUserPersona, "")
	// 	if err != nil {
	// 		return fmt.Errorf("could not load active user persona '%s': %s", p.Api.Tympan.Configuration.ActiveUserPersona, err)
	// 	}

	// 	activeSkirmish, err := p.Api.GetActiveSkirmish(activeUserPersona, "")
	// 	if err != nil {
	// 		return fmt.Errorf("could not load active skirmish '%s': %s", activeUserPersona.Settings.ActiveSkirmish, err)
	// 	}
	// 	fmt.Printf("Congrats! You're playing a skirmish! This one:\n%+v", activeSkirmish)
	// 	// wire up so you _actually_ play.
	// }

	// Turn off for now
	// group, err := prompts.NewGroup(p.Api.CachedProfiles, p.Api.CachedTraits)
	// if err != nil {
	// 	return err
	// }

	// switch viper.GetString("format") {
	// case "json":
	// 	fmt.Print(group.JSON())
	// default:
	// 	groupOutput := strings.Builder{}
	// 	groupOutput.WriteString("| Name | Profile | Melee | Missile | Move | FS | R | T | Traits |\n")
	// 	groupOutput.WriteString("| ---- | ------- | ----- | ------- | ---- | -- | - | - | ------ |\n")
	// 	groupOutput.WriteString(group.MarkdownTableEntry())
	// 	output, err := p.Dossier.Render(groupOutput.String())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Print(output)
	// }
	return nil
}
