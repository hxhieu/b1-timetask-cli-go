//go:build linux

package cmd

import (
	"github.com/hxhieu/b1-timetask-cli-go/console"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	if !ctx.Experimental {
		console.ErrorLn("Native coming soon (TM)")
		return nil
	}

	// DEBUG MODE ONLY

	return nil
}
