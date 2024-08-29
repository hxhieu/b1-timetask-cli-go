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

func (t *ToolbarLabel) ToolbarObject() fyne.CanvasObject {
	return t.label
}

func NewToolbarLabel(text string) *ToolbarLabel {
	b := binding.NewString()
	b.Set(text)
	l := widget.NewLabelWithData(b)
	return &ToolbarLabel{
		textBinding: &b,
		label:       l,
	}
}

func (t *ToolbarLabel) SetText(text string) {
	(*t.textBinding).Set(text)
}
