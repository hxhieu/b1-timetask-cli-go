//go:build linux || android

package cmd

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	// if !ctx.Experimental {
	// 	console.ErrorLn("Native coming soon (TM)")
	// 	return nil
	// }

	// DEBUG MODE ONLY

	a := app.New()
	w := a.NewWindow("test")
	defer w.ShowAndRun()
	w.SetContent(
		widget.NewButton("Click me!!!", nil),
	)

	return nil
}
