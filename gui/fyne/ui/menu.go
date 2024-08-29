package ui

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	uiWidget "github.com/hxhieu/b1-timetask-cli-go/gui/fyne/widget"
)

func NewMenu(appState *AppState) *widget.Toolbar {
	userLabel := uiWidget.NewToolbarLabel(&appState.User)
	spinner := uiWidget.NewToolbarActivity()
	actionButton := widget.NewToolbarAction(theme.LoginIcon(), func() {
		if user, _ := appState.User.Get(); user == "Login" {
			tokenEntry := widget.NewEntry()
			dialog.ShowForm(
				"Login with token",
				"Login",
				"Cancel",
				[]*widget.FormItem{
					{
						Text:   "Token",
						Widget: tokenEntry,
					},
				},
				func(b bool) {
					fmt.Printf("%s\n", tokenEntry.Text)
				},
				appState.MainWindow,
			)
		} else {

		}
	})
	actionButton.ToolbarObject().Hide()

	// Check user
	go func(l *uiWidget.ToolbarActivity, state *AppState) {
		if user, err := state.Common.InitUser(); err == nil {
			state.User.Set(*user)
			actionButton.Icon = theme.AccountIcon()
		} else {
			state.User.Set("Login")
			fmt.Println(err.Error())
			dialog.ShowError(err, state.MainWindow)
		}
		actionButton.ToolbarObject().Show()
		l.Stop()
	}(spinner, appState)

	toolbar := widget.NewToolbar(
		actionButton,
		userLabel,
		widget.NewToolbarSpacer(),
		spinner,
	)

	return toolbar
}
