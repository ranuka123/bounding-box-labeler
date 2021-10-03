package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type interactiveRaster struct {
	widget.BaseWidget

	label *labeler
	min   fyne.Size
	img   *canvas.Raster

	dragStart bool
	startX    int
	startY    int
	endX      int
	endY      int

	tappedLoc []int
}

func (r *interactiveRaster) SetMinSize(size fyne.Size) {
	pixWidth, _ := r.locationForPosition(fyne.NewPos(size.Width, size.Height))
	scale := float32(1.0)
	c := fyne.CurrentApp().Driver().CanvasForObject(r.img)
	if c != nil {
		scale = c.Scale()
	}

	texScale := float32(pixWidth) / size.Width * float32(r.label.zoom) / scale
	size = fyne.NewSize(size.Width/texScale, size.Height/texScale)
	r.min = size
	r.Resize(size)
}

func (r *interactiveRaster) MinSize() fyne.Size {
	return r.min
}

func (r *interactiveRaster) CreateRenderer() fyne.WidgetRenderer {
	return &rasterWidgetRender{raster: r, bg: canvas.NewRasterWithPixels(bgPattern)}
}

func (r *interactiveRaster) Dragged(event *fyne.DragEvent) {
	if r.label.img == nil {
		return
	}

	x, y := r.locationForPosition(event.PointEvent.Position)
	if x >= r.label.img.Bounds().Dx() || y >= r.label.img.Bounds().Dy() {
		return
	}

	if r.dragStart {
		r.endX = x
		r.endY = y
	} else {
		r.startX = x
		r.startY = y
		r.dragStart = true
	}
}

func (r *interactiveRaster) DragEnd() {
	if r.label.img == nil {
		return
	}
	r.dragStart = false
	r.label.AddBox(r.startX, r.startY, r.endX, r.endY)
}

func (r *interactiveRaster) locationForPosition(pos fyne.Position) (int, int) {
	c := fyne.CurrentApp().Driver().CanvasForObject(r.img)
	x, y := int(pos.X), int(pos.Y)
	if c != nil {
		x, y = c.PixelCoordinateForPosition(pos)
	}

	return x / r.label.zoom, y / r.label.zoom
}

func newInteractiveRaster(label *labeler) *interactiveRaster {
	r := &interactiveRaster{img: canvas.NewRaster(label.draw), label: label}
	r.ExtendBaseWidget(r)
	return r
}

type rasterWidgetRender struct {
	raster *interactiveRaster
	bg     *canvas.Raster
}

func bgPattern(x, y, _, _ int) color.Color {
	const boxSize = 25

	if (x/boxSize)%2 == (y/boxSize)%2 {
		return color.Gray{Y: 58}
	}

	return color.Gray{Y: 84}
}

func (r *rasterWidgetRender) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.raster.img.Resize(size)
}

func (r *rasterWidgetRender) MinSize() fyne.Size {
	return r.MinSize()
}

func (r *rasterWidgetRender) Refresh() {
	canvas.Refresh(r.raster)
}

func (r *rasterWidgetRender) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *rasterWidgetRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.raster.img}
}

func (r *rasterWidgetRender) Destroy() {
}
