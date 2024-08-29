package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ToolbarActivity struct {
	activity widget.Activity
}

// Implements ToolbarObject
func (t *ToolbarActivity) ToolbarObject() fyne.CanvasObject {
	return &t.activity
}

func NewToolbarActivity() *ToolbarActivity {
	t := ToolbarActivity{
		activity: *widget.NewActivity(),
	}
	t.activity.Start()
	return &t
}

func (t *ToolbarActivity) Start() {
	t.activity.Show()
}

func (t *ToolbarActivity) Stop() {
	t.activity.Hide()
}
