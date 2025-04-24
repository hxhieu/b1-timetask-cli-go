package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewTaskTable(appState *AppState) (*widget.Table, error) {
	tasks, err := appState.Common.FetchTaskInputs()
	if err != nil {
		return nil, err
	}

	list := widget.NewTable(
		func() (int, int) {
			return len(tasks), 11 // Task input has 11 fields
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("No content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			// o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	return list, nil
}
