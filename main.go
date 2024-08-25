package main

import (
	"embed"

	"github.com/alecthomas/kong"
	"github.com/hxhieu/b1-timetask-cli-go/cmd"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(
		&cli,
		kong.Description("A CLI tool to semi-automate the creation of time task."),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			NoAppSummary: true,
		}),
	)
	err := ctx.Run(
		// Bind the context with the flags
		cmd.CLIContext{
			Debug:        cli.Debug,
			Force:        cli.Force,
			Experimental: cli.Experimental,
			GuiAssets:    &assets,
		},
	)
	ctx.FatalIfErrorf(err)
}
