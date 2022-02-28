package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCommand struct {
	ConfigFile            string
	LogLevel              string
	Format                string
	Api                   *flfa.Api
	DefaultConfigFileName string
}

type RootCommander interface {
	CreateCommand() *cobra.Command
	InitLogger()
	InitConfig()
}

func (r *RootCommand) CreateCommand() *cobra.Command {

	tmp := &cobra.Command{
		Use:              "flfa",
		Short:            "flfa - Flagrant Factions",
		Long:             `Flagrant Factions (flfa) - play Flagrant Factions from your terminal`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	tmp.PersistentFlags().StringVar(&r.ConfigFile, "config", "", fmt.Sprintf("config file (default is $HOME/.config/%s)", r.DefaultConfigFileName))

	tmp.PersistentFlags().StringVar(&r.LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	err := tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	tmp.PersistentFlags().StringVar(&r.Format, "format", "human", "Output format: human readable or JSON")
	err = tmp.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"human", "json"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)
	viper.BindPFlag("format", tmp.PersistentFlags().Lookup("format"))

	return tmp
}

func (r *RootCommand) InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lvl, err := zerolog.ParseLevel(r.LogLevel)
	if err != nil {
		panic("Could not initialize zerolog")
	}

	zerolog.SetGlobalLevel(lvl)

	if lvl == zerolog.InfoLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	log.Trace().Msg("Initialized zerolog")
}

func (r *RootCommand) InitConfig() {
	if r.ConfigFile != "" {
		viper.SetConfigFile(r.ConfigFile)
	} else {
		home, _ := homedir.Dir()
		r.ConfigFile = r.DefaultConfigFileName
		viper.SetConfigName(r.ConfigFile)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		ConfigFileDirectoryPath := filepath.Join(home, ".config")
		viper.AddConfigPath(ConfigFileDirectoryPath)

		if _, err := os.Stat(ConfigFileDirectoryPath); os.IsNotExist(err) {
			log.Trace().Msgf("%s does not exist, creating", ConfigFileDirectoryPath)
			if err := os.MkdirAll(ConfigFileDirectoryPath, 0750); err != nil {
				log.Error().Msgf("failed to create dir %s: %s", ConfigFileDirectoryPath, err)
			}
		}

		ConfigFilePath := filepath.Join(ConfigFileDirectoryPath, r.ConfigFile)

		if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
			_, err := os.Create(filepath.Clean(ConfigFilePath))
			if err != nil {
				log.Error().Msgf("failed to initialise %s: %s", ConfigFilePath, err)
			}
			r.Api.GenerateDefaultConfig()
			viper.WriteConfig()
		}
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err == nil {
		log.Trace().Msgf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Trace().Msgf("Could not load config file '%s' because: %s", r.ConfigFile, err)
		log.Trace().Msg("Falling back on default config")
		r.Api.GenerateDefaultConfig()
	}

	if err := r.Api.LoadConfig(); err != nil {
		log.Warn().Msgf("Error setting running config: %s", err)
	}
}

// Returns the cobra command called, e.g. new or install
// and also the fully formatted command as passed with arguments/flags.
// Idea borrowed from carolynvs/porter:
// https://github.com/carolynvs/porter/blob/ccca10a63627e328616c1006600153da8411a438/cmd/porter/main.go
func GetCalledCommand(cmd *cobra.Command) (string, string) {
	if len(os.Args) < 2 {
		return "", ""
	}

	calledCommandStr := os.Args[1]

	// Also figure out the full called command from the CLI
	// Is there anything sensitive you could leak here? 🤔
	calledCommandArgs := strings.Join(os.Args[1:], " ")

	return calledCommandStr, calledCommandArgs
}
