package cmd

import "github.com/hxhieu/b1-timetask-cli-go/console"

type guiCmd struct {
}

func (c *guiCmd) Run(ctx CLIContext) error {
	if !ctx.Debug {
		console.ErrorLn("Coming soon (TM)")
	}

	// DEBUG MODE ONLY

	return nil
}
