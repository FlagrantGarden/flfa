package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCommand contains all the necessary information to initialize a Tympan application
type RootCommand[Config tympan.Configurable] struct {
	ConfigFile string
	LogLevel   string
	Format     string
	Tympan     tympan.Tympan[Config]
}

// The Version command is reimplementable and must return a valid cobra Command, initialize logging, and intialize
// application configuration.
type RootCommander interface {
	CreateCommand() *cobra.Command
	InitLogger()
	InitConfig()
}

// CreateCommand instatiates the root cobra command of a Tympan app. All other commands are inherited from this one. The
// root command is where overridable application-wide behavior is implemented; it allows a user to pass the --config
// flag to specify their own configuration file instead of the default one, the --log-level flag to tune the verbosity
// of the application, and the --format flag to choose between human readable and JSON output.
//
// It also uses the Tympan Metadata to fill out the root commands use, short name, and long name.
func (r *RootCommand[Config]) CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:              r.Tympan.Metadata.Name,
		Short:            fmt.Sprintf("%s (%s)", r.Tympan.Metadata.Name, r.Tympan.Metadata.DisplayName),
		Long:             fmt.Sprintf("%s (%s) - %s", r.Tympan.Metadata.DisplayName, r.Tympan.Metadata.Name, r.Tympan.Metadata.Description),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	defaultConfigFilePath := r.Tympan.DefaultConfigFile()
	tmp.PersistentFlags().StringVar(&r.ConfigFile, "config", defaultConfigFilePath, "configuration file path")

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

// Configures the zerolog logging for the application, enabling it to be propagated to child commands.
func (r *RootCommand[Config]) InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}

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

// Initializes the applications configuration (including both the shared configuration from Tympan and the application
// developer's specified configuration), writing it to disk if it doesn't exist and loading it if it does.
func (r *RootCommand[Config]) InitConfig() {
	err := r.Tympan.InitializeConfig()
	if err != nil {
		log.Error().Msgf("unable to initialize configuration: %s", err)
	}
	log.Logger.Trace().Msgf("Using cache path: %s", r.Tympan.Configuration.GetFolderPath("cache"))
	log.Logger.Trace().Msgf("Using application path: %s", r.Tympan.Configuration.GetFolderPath("application"))
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
