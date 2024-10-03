package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hxhieu/b1-timetask-cli-go/gui/fyne/ui"
)

// Initial state
var appState *ui.AppState

func NewNativeGui(title string) {
	appState = ui.NewAppState()
	w := appState.App.NewWindow(title)
	w.SetMaster()
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(1024, 768))
	defer w.ShowAndRun()
	appState.MainWindow = w

	toolbar := ui.NewMenu(appState)

	content := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("Content"))

	w.SetContent(content)
}
