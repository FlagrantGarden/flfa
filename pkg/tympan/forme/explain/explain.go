package explain

import (
	"fmt"

	"github.com/FlagrantGarden/flfa/pkg/tympan/dossier"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ExplainCommand contains all the necessary information to display non-reference documentation at the terminal
type ExplainCommand struct {
	Dossier    *dossier.Dossier
	DocsUri    string
	Online     bool
	ListTopics bool
	Tag        string
	Category   string
	Topic      string
}

// Cobra expects a string that explains the usage of a command and recommends the following syntax:
//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
//   ... indicates that you can specify multiple values for the previous argument.
//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
//       argument to the right of the separator. You cannot use both arguments in a single use of the command.
//   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
//       optional, they are enclosed in brackets ([ ]).
var Use = "explain [topic] [--list]"

// This string is shown in the help output of the root command.
var Short = "Present documentation about topics"

// this string is shown in the help output of this command.
var Long = heredoc.Doc(`
	Present documentation about various topics. You can list available topics,
	filter the list of available topics, render a single topic to the terminal,
	and automatically navigate to the webpage for a document. By default, it
	lists all available topics in a table with their topic, description category,
	and tags.
`)

// Create command instantiates a cobra command that can be added to an application to render non-reference documentation
// at the terminal or send the user to the appropriate web page.
func (e *ExplainCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               Use,
		Short:             Short,
		Long:              Long,
		Args:              e.validateArgCount,
		ValidArgsFunction: e.completeTitle,
		PreRun:            e.preExecute,
		RunE:              e.execute,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVarP(&e.ListTopics, "list", "l", false, "list available topics")

	cmd.Flags().StringVarP(&e.Tag, "tag", "t", "", "filter available topics by tag")
	err := cmd.RegisterFlagCompletionFunc("tag", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if e.Dossier.Documents == nil {
			e.preExecute(cmd, args)
		}
		return e.Dossier.ListTags(e.Dossier.Documents), cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	cmd.Flags().StringVarP(&e.Category, "category", "c", "", "filter available topics by category")
	err = cmd.RegisterFlagCompletionFunc("category", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if e.Dossier.Documents == nil {
			e.preExecute(cmd, args)
		}
		return e.Dossier.ListCategories(e.Dossier.Documents), cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	if e.DocsUri != "" {
		cmd.Flags().BoolVar(&e.Online, "online", false, "view documentation online")
	}

	return cmd
}

// preExecute searches for documents, parsing and cachine any it discovers for use in the rest of the command.
func (e *ExplainCommand) preExecute(cmd *cobra.Command, args []string) {
	e.Dossier.CacheDocuments("content")
}

// validateArgCount ensures that the user is passing a valid group of flags; if the user does not specify a topic *or*
// the --list flag, it treats the command as if the list flag was passed. If the user specifies a topic and a filter
// flag, it errors and reminds the user to specify *either* a topic or filter flags. If the user specifies more than one
// topic, it errors and reminds the user to specify only one topic at a time.
func (e *ExplainCommand) validateArgCount(cmd *cobra.Command, args []string) error {
	// show available topics if user runs `pct explain`
	if len(args) == 0 && !e.ListTopics {
		e.ListTopics = true
	}

	if len(args) == 1 {
		if e.Category != "" || e.Tag != "" {
			return fmt.Errorf("specify a topic *or* search by tag/category")
		}
		e.Topic = args[0]
	} else if len(args) > 1 {
		return fmt.Errorf("specify only one topic to explain")
	}

	return nil
}

// completeTitle is used in shell completion to recommend valid titles that match what the user has already typed for
// the argument.
func (e *ExplainCommand) completeTitle(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if e.Dossier.Documents == nil {
		e.preExecute(cmd, args)
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return e.Dossier.CompleteShortTitle(toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

// execute handles the discovery and displaying of documentation. If the --list flag is passed or no topic is specified,
// it returns the available documents (possibly filtered by tag and/or category). If only one document matches a given
// search filter combination, that document is rendered. If a single topic is passed, execute renders that document in
// the terminal unless the --online flag was passed in which case it instead opens the browser page for that document
// on the project website. If the output format is set to JSON, it returns a JSON object instead of rendering to the
// terminal.
func (e *ExplainCommand) execute(cmd *cobra.Command, args []string) error {
	docs := e.Dossier.Documents
	if e.ListTopics {
		format := viper.GetString("format")
		if format == "" {
			format = "table"
		}
		if e.Category != "" {
			docs = dossier.FilterByCategory(e.Category, docs)
		}
		if e.Tag != "" {
			docs = dossier.FilterByTag(e.Tag, docs)
		}
		// If there's only one match, should we render the matching doc?
		if len(docs) == 1 {
			output, err := e.Dossier.RenderDocument(docs[0])
			if err != nil {
				return err
			}
			fmt.Print(output)
		} else {
			e.Dossier.FormatFrontMatter(format, docs)
		}
	} else if e.Topic != "" {
		doc, err := e.Dossier.SelectDocument(e.Topic)
		if err != nil {
			return err
		}
		output, err := e.Dossier.RenderDocument(doc)
		if err != nil {
			return err
		}
		fmt.Print(output)
		// If --online, open in browser and do not display
		// Should we have a --scroll mode?
	}
	return nil
}
