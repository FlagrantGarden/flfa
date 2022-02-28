package version

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type VersionCommand struct {
	Version   string
	BuildDate string
	Commit    string
	SourceUrl string
}

type VersionCommander interface {
	CreateCommand() *cobra.Command
}

func (v *VersionCommand) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, Format(v.Version, v.BuildDate, v.Commit, v.SourceUrl))
		},
	}

	return cmd
}

func Format(version string, buildDate string, commit string, sourceUrl string) string {
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

	return fmt.Sprintf("flfa %s %s %s\n\n%s", version, commit, dateStr, changelogURL(version, sourceUrl))
}

func changelogURL(version string, sourceUrl string) string {
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", sourceUrl)
	}

	url := fmt.Sprintf("%s/releases/tag/%s", sourceUrl, strings.TrimPrefix(version, "v"))
	return url
}
