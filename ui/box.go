package ui

import (
	"bbox_labeler/image"
	"fmt"
)

var turquoise []uint8 = []uint8{0xaf, 0xee, 0xee, 0xff}

func (l *labeler) SetBoxList(boxes []image.Box) {

	boxList := []string{}

	for _, box := range boxes {
		boxList = append(boxList, fmt.Sprintf("%d %d %d %d", box[0], box[1], box[2], box[3]))
	}
	l.boxList.Set(boxList)

}

func (l *labeler) AddBox(x, y, x1, y1 int) {

	img := image.Get(l.index)
	img.AddBox(x, y, x1, y1)

	l.DrawBox(x, y, x1, y1)
	l.SetBoxList(img.Boxes)
}

func (l *labeler) DrawBox(x, y, x1, y1 int) {

	if x < x1 {
		l.DrawLine(x, y, x1, y)
		l.DrawLine(x, y1, x1, y1)
	} else {
		l.DrawLine(x1, y, x, y)
		l.DrawLine(x1, y1, x, y1)
	}

	if y < y1 {
		l.DrawLine(x, y, x, y1)
		l.DrawLine(x1, y, x1, y1)
	} else {
		l.DrawLine(x, y1, x, y)
		l.DrawLine(x1, y1, x1, y)
	}

	l.renderCache()
}

func (l *labeler) DrawLine(x, y, x1, y1 int) {

	for i := x; i <= x1; i++ {
		for j := y; j <= y1; j++ {
			l.SetPixelColor(i, j, turquoise)
		}
	}

}
