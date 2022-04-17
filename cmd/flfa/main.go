/*
Copyright © 2022 Flagrant Garden flagrant.garden@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"

	"github.com/FlagrantGarden/flfa/cmd/flfa/editor"
	"github.com/FlagrantGarden/flfa/cmd/flfa/play"
	"github.com/FlagrantGarden/flfa/docs"
	"github.com/FlagrantGarden/flfa/emfs"
	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/FlagrantGarden/flfa/pkg/tympan/dossier"
	"github.com/FlagrantGarden/flfa/pkg/tympan/forme/explain"
	"github.com/FlagrantGarden/flfa/pkg/tympan/forme/root"
	"github.com/FlagrantGarden/flfa/pkg/tympan/forme/version"
	"github.com/FlagrantGarden/flfa/pkg/tympan/telemetry"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	app_version       = "dev"
	commit            = "none"
	date              = "unknown"
	sourceUrl         = "https://github.com/FlagrantGarden/flfa"
	honeycomb_api_key = "not_set"
	honeycomb_dataset = "not_set"
)

func main() {
	// Telemetry must be initialized before anything else;
	// If the telemetry build tag was not passed, this is all null ops
	ctx, traceProvider, parentSpan := telemetry.Start(context.Background(), honeycomb_api_key, honeycomb_dataset, "flfa-cli", "flfa")

	// Create FLFA context
	fs := afero.NewOsFs() // use the real file system
	mfs := emfs.GetEmbeddedModulesFS()
	api := &flfa.Api{
		EMFS: &mfs,
		Tympan: &tympan.Tympan[*flfa.Configuration]{
			AFS:           &afero.Afero{Fs: fs},
			Configuration: &flfa.Configuration{},
			Metadata: tympan.Metadata{
				Name:           "flfa",
				DisplayName:    "Flagrant Factions",
				Description:    "Play Flagrant Factions from your terminal",
				FolderName:     "FlagrantFactions",
				ConfigFileName: "config",
				SourceUrl:      "https://github.com/FlagrantGarden/flfa",
				ProjectUrl:     "https://flagrant.garden/games/factions",
			},
		},
	}
	// Setup the root command (flfa)
	root_cmder := &root.RootCommand[*flfa.Configuration]{
		Tympan: *api.Tympan,
	}
	root_cmd := root_cmder.CreateCommand()
	version_string := version.Format(api.Tympan.Metadata.Name, app_version, date, commit, api.Tympan.Metadata.SourceUrl)
	root_cmd.Version = version_string
	root_cmd.SetVersionTemplate(version_string)

	// Get the command called and its arguments;
	// The arguments are only necessary if we want to
	// hand them off as an attribute to the parent span:
	// do we? Otherwise we just need the calledCommand
	calledCommand, _ := root.GetCalledCommand(root_cmd)

	// flfa version
	version_cmder := &version.VersionCommand{
		Version:   app_version,
		BuildDate: date,
		Commit:    commit,
		Metadata:  api.Tympan.Metadata,
	}
	root_cmd.AddCommand(version_cmder.CreateCommand())

	// flfa explain
	efs := docs.GetDocsFS()
	explainer := explain.ExplainCommand{
		Dossier: &dossier.Dossier{
			AFS: api.Tympan.AFS,
			EFS: &efs,
		},
	}
	root_cmd.AddCommand(explainer.CreateCommand())

	// flfa play
	play_cmder := play.PlayCommand{
		Api:     api,
		Dossier: &dossier.Dossier{},
	}
	play_cmd := play_cmder.CreateCommand()
	root_cmd.AddCommand(play_cmd)

	// flfa editor

	editor_cmder := editor.EditorCommand{
		Api: api,
	}
	editor_cmd := editor_cmder.CreateCommand()
	root_cmd.AddCommand(editor_cmd)

	// initialize
	cobra.OnInitialize(root_cmder.InitLogger, root_cmder.InitConfig)

	// instrument & execute called command
	ctx, childSpan := telemetry.NewSpan(ctx, calledCommand)
	err := root_cmd.ExecuteContext(ctx)
	telemetry.RecordSpanError(childSpan, err)
	telemetry.EndSpan(childSpan)

	// Send all events
	telemetry.ShutDown(ctx, traceProvider, parentSpan)

	// handle exiting with/out errors
	cobra.CheckErr(err)
}
