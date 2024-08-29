package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ToolbarLabel struct {
	textBinding *binding.String
	label       *widget.Label
}

// Implements ToolbarObject
func (t *ToolbarLabel) ToolbarObject() fyne.CanvasObject {
	return t.label
}

func NewToolbarLabel(text *binding.String) *ToolbarLabel {
	l := widget.NewLabelWithData(*text)
	return &ToolbarLabel{
		textBinding: text,
		label:       l,
	}
}

func (t *ToolbarLabel) SetText(text string) {
	(*t.textBinding).Set(text)
}
