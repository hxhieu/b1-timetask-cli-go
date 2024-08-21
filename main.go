package main

import (
	"github.com/alecthomas/kong"
	"github.com/hxhieu/b1-timetask-cli-go/cmd"
)

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(
		&cli,
		kong.Name("/path/to/the/cli"),
		kong.Description("A CLI tool to semi-automate the creation of time task."),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			NoAppSummary: true,
		}),
	)
	err := ctx.Run(cmd.CLIContext{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
