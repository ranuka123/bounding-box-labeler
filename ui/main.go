package ui

import (
	"log"

	images "bbox_labeler/image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (l *labeler) Update() {
	im := images.Get(l.index)
	l.LoadImage(im)
	l.fileName.Set(im.Name)
	l.SetBoxList(im.Boxes)

	for _, box := range im.Boxes {
		l.DrawBox(box[0], box[1], box[2], box[3])
	}
}

func (l *labeler) BuildUI(w fyne.Window) {
	l.win = w

	/*
	 Inputs
	*/
	directoryEntry := widget.NewEntry()

	/*
	 Buttons
	*/
	prevBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if l.index > 0 {
			l.index -= 1
			l.Update()
		}
	})
	nextBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		if l.index < images.Size() {
			l.index += 1
			l.Update()
		}
	})
	resetBtn := widget.NewButtonWithIcon("Clear Boxes", theme.DeleteIcon(), func() {
		images.Get(l.index).ClearBoxes()
		l.Update()
	})
	saveAllBtn := widget.NewButtonWithIcon("Write All", theme.DocumentSaveIcon(), func() {
		if err := images.WriteBoxes(); err != nil {
			log.Println(err)
		}
	})
	loadDirectoryBtn := widget.NewButtonWithIcon("Load Directory", theme.DocumentIcon(), func() {
		err := images.LoadImages(directoryEntry.Text)
		if err != nil {
			log.Println(err)
			return
		}
		if err := images.LoadBoxes(directoryEntry.Text); err != nil {
			log.Println(err)
		}

		l.Update()
	})

	/*
	 Labels
	*/
	l.fileName = binding.NewString()
	fileLabel := widget.NewLabelWithData(l.fileName)

	/*
	 Tables
	*/
	l.boxList = binding.NewStringList()
	list := widget.NewListWithData(l.boxList,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	/*
	 Layout
	*/
	leftColHeader := container.NewVBox(
		container.New(layout.NewGridLayout(2), prevBtn, nextBtn),
		resetBtn,
		fileLabel)
	leftCol := container.New(layout.NewBorderLayout(leftColHeader, nil, nil, nil), leftColHeader, list)
	content := container.NewHSplit(leftCol, container.NewScroll(l.drawSurface))
	content.Offset = 0.2

	header := container.NewHSplit(directoryEntry, container.NewHBox(loadDirectoryBtn, layout.NewSpacer(), saveAllBtn))

	w.SetContent(container.New(layout.NewBorderLayout(header, nil, nil, nil), header, content))
}
