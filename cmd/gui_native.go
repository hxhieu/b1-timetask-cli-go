//go:build linux || android

package cmd

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	myWidget "github.com/hxhieu/b1-timetask-cli-go/gui/fyne/widget"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	// if !ctx.Experimental {
	// 	console.ErrorLn("Native coming soon (TM)")
	// 	return nil
	// }

	// DEBUG MODE ONLY

	a := app.New()
	w := a.NewWindow(APP_TITLE)
	defer w.ShowAndRun()

	userLabel := myWidget.NewToolbarLabel("11")

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		userLabel,
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			log.Println("Display help")
		}),
	)

	toolbar.Items[0].ToolbarObject().Hide()
	content := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("Content"))

	w.SetMaster()
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(1024, 768))
	w.SetContent(content)

	return nil
}
