package ui

import (
	im "bbox_labeler/image"
	"image"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/draw"
)

type labeler struct {
	drawSurface             *interactiveRaster
	status                  *widget.Label
	cache                   *image.RGBA
	cacheWidth, cacheHeight int
	fgPreview               *canvas.Rectangle

	uri  string
	img  *image.RGBA
	zoom int
	fg   color.Color

	win        fyne.Window
	recentMenu *fyne.Menu
	index      int
	boxes      *[]string
	boxList    binding.StringList
	fileName   binding.String
}

func colorToBytes(col color.Color) []uint8 {
	r, g, b, a := col.RGBA()
	return []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func (l *labeler) SetPixelColor(x, y int, rgba []uint8) {
	i := (y*l.img.Bounds().Dx() + x) * 4
	l.img.Pix[i] = rgba[0]
	l.img.Pix[i+1] = rgba[1]
	l.img.Pix[i+2] = rgba[2]
	l.img.Pix[i+3] = rgba[3]
}

func (l *labeler) buildUI() fyne.CanvasObject {
	return container.NewScroll(l.drawSurface)
}

func (l *labeler) draw(w, h int) image.Image {
	if l.cacheWidth == 0 || l.cacheHeight == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	if w > l.cacheWidth || h > l.cacheHeight {
		bigger := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(bigger, l.cache.Bounds(), l.cache, image.Point{}, draw.Over)
		return bigger
	}

	return l.cache
}

func (l *labeler) updateSizes() {
	if l.img == nil {
		return
	}
	l.cacheWidth = l.img.Bounds().Dx() * l.zoom
	l.cacheHeight = l.img.Bounds().Dy() * l.zoom

	c := fyne.CurrentApp().Driver().CanvasForObject(l.status)
	scale := float32(1.0)
	if c != nil {
		scale = c.Scale()
	}
	l.drawSurface.SetMinSize(fyne.NewSize(
		float32(l.cacheWidth)/scale,
		float32(l.cacheHeight)/scale))

	l.renderCache()
}

func (l *labeler) pixAt(x, y int) []uint8 {
	ix := x / l.zoom
	iy := y / l.zoom

	if ix >= l.img.Bounds().Dx() || iy >= l.img.Bounds().Dy() {
		return []uint8{0, 0, 0, 0}
	}

	return colorToBytes(l.img.At(ix, iy))
}

func (l *labeler) renderCache() {
	l.cache = image.NewRGBA(image.Rect(0, 0, l.cacheWidth, l.cacheHeight))
	for y := 0; y < l.cacheHeight; y++ {
		for x := 0; x < l.cacheWidth; x++ {
			i := (y*l.cacheWidth + x) * 4
			col := l.pixAt(x, y)
			l.cache.Pix[i] = col[0]
			l.cache.Pix[i+1] = col[1]
			l.cache.Pix[i+2] = col[2]
			l.cache.Pix[i+3] = col[3]
		}
	}

	l.drawSurface.Refresh()
}

func fixEncoding(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)
	return newImg
}

func (l *labeler) LoadImage(i *im.Image) {

	img, err := i.Image()
	if err != nil {
		log.Println(err)
		return
	}

	l.img = fixEncoding(*img)
	l.updateSizes()
}

func NewEditor() *labeler {

	fgCol := color.Black
	label := &labeler{zoom: 1, fg: fgCol, fgPreview: canvas.NewRectangle(fgCol)}
	label.drawSurface = newInteractiveRaster(label)

	return label
}
