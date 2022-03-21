package version

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/FlagrantGarden/flfa/pkg/tympan"

	"github.com/spf13/cobra"
)

// VersionCommand contains all the necessary information to determine and display a Tympan application's version.
type VersionCommand struct {
	Version   string
	BuildDate string
	Commit    string
	tympan.Metadata
}

// The Version command is reimplementable and must return a valid cobra Command.
type VersionCommander interface {
	CreateCommand() *cobra.Command
}

// CreateCommand instantiates a cobra command to display the version of the Tympan application.
func (v *VersionCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, Format(v.Name, v.Version, v.BuildDate, v.Commit, v.SourceUrl))
		},
	}

	return cmd
}

// Format is used to display version information about a Tympan application. For an application named "myapp" with
// a version of "1.2.3", a build date of "2022/01/31", a commit hash of "f1976f1", and a source url of
// https://github.com/myuser/myapp, it will return:
//     `
//     myapp 1.2.3 f1976f1 2022/01/31
//     https://github.com/myuser/myapp/releases/tag/1.2.3
//     `
func Format(name string, version string, buildDate string, commit string, sourceUrl string) string {
	version = strings.TrimSpace(strings.TrimPrefix(version, "v"))

	var dateStr string
	if buildDate != "" {
		t, _ := time.Parse(time.RFC3339, buildDate)
		dateStr = t.Format("2006/01/02")
	}

	if commit != "" && len(commit) > 7 {
		length := len(commit) - 7
		commit = strings.TrimSpace(commit[:len(commit)-length])
	}

	return fmt.Sprintf("%s %s %s %s\n\n%s", name, version, commit, dateStr, githubReleaseTagURL(version, sourceUrl))
}

// githubReleaseTagURL is a helper function to determine the URL to the appropriate release of the application on GitHub
func githubReleaseTagURL(version string, sourceUrl string) string {
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", sourceUrl)
	}

	url := fmt.Sprintf("%s/releases/tag/%s", sourceUrl, strings.TrimPrefix(version, "v"))
	return url
}
