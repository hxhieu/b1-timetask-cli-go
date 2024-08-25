//go:build linux

package cmd

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/hxhieu/b1-timetask-cli-go/console"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	if !ctx.Experimental {
		console.ErrorLn("Native coming soon (TM)")
		return nil
	}

	// DEBUG MODE ONLY

	a := app.New()
	w := a.NewWindow("B1 TimeTask CLI GUI")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))
	w.Resize(fyne.NewSize(1024, 768))

	w.ShowAndRun()

	return nil
}
