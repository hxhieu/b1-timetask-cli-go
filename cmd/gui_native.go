//go:build linux || android

package cmd

import (
	mygui "github.com/hxhieu/b1-timetask-cli-go/gui/fyne"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	mygui.NewNativeGui(APP_TITLE)
	return nil
}
