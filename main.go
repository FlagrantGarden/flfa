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

	cmdExplain "github.com/FlagrantGarden/flfa/cmd/explain"
	cmdPlay "github.com/FlagrantGarden/flfa/cmd/play"
	cmdRoot "github.com/FlagrantGarden/flfa/cmd/root"
	cmdVersion "github.com/FlagrantGarden/flfa/cmd/version"
	"github.com/FlagrantGarden/flfa/pkg/flfa"
	"github.com/FlagrantGarden/flfa/pkg/telemetry"
	"github.com/FlagrantGarden/flfa/pkg/tympan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	version           = "dev"
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
	api := &flfa.Api{
		Tympan: &tympan.Tympan{
			AFS:  &afero.Afero{Fs: fs},
			IOFS: &afero.IOFS{Fs: fs},
		},
	}
	// Setup the root command (flfa)
	root_cmder := &cmdRoot.RootCommand{
		Api:                   api,
		DefaultConfigFileName: ".flagrant-factions.yaml",
	}
	root_cmd := root_cmder.CreateCommand()
	version_string := cmdVersion.Format(version, date, commit, sourceUrl)
	root_cmd.Version = version_string
	root_cmd.SetVersionTemplate(version_string)
	// Get the command called and its arguments;
	// The arguments are only necessary if we want to
	// hand them off as an attribute to the parent span:
	// do we? Otherwise we just need the calledCommand
	calledCommand, _ := cmdRoot.GetCalledCommand(root_cmd)

	// flfa version
	version_cmder := &cmdVersion.VersionCommand{
		Version:   version,
		BuildDate: date,
		Commit:    commit,
		SourceUrl: sourceUrl,
	}
	version_cmd := version_cmder.CreateCommand()
	root_cmd.AddCommand(version_cmd)

	// flfa explain
	root_cmd.AddCommand(cmdExplain.CreateCommand())

	// flfa play
	play_cmder := cmdPlay.PlayCommand{
		Api: api,
	}
	play_cmd := play_cmder.CreateCommand()
	root_cmd.AddCommand(play_cmd)

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
