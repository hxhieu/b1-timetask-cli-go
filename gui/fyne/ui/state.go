package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"github.com/hxhieu/b1-timetask-cli-go/gui"
)

type AppState struct {
	Common     gui.App
	User       binding.String
	App        fyne.App
	MainWindow fyne.Window
}

func NewAppState() *AppState {
	return &AppState{
		Common: gui.App{},
		User:   binding.NewString(),
		App:    app.New(),
	}
}
