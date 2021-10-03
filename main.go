package main

import (
	"bbox_labeler/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	a := app.NewWithID("bounding box labeler")
	w := a.NewWindow("BBox Labeler")

	e := ui.NewEditor()
	e.BuildUI(w)

	w.Resize(fyne.NewSize(800, 800))
	w.ShowAndRun()

}
