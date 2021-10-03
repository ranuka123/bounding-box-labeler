package image

import (
	"fmt"
	im "image"
	"os"
)

type Box [4]int

type Image struct {
	Name  string
	Boxes []Box
}

func (i *Image) AddBox(x, y, x1, y1 int) {
	i.Boxes = append(i.Boxes, Box{x, y, x1, y1})
}

func (i *Image) ClearBoxes() {
	i.Boxes = []Box{}
}

func (i *Image) Image() (*im.Image, error) {
	f, err := os.Open(fmt.Sprintf("%s%s", imagePath, i.Name))
	if err != nil {
		panic(err)
	}

	defer f.Close()

	img, _, err := im.Decode(f)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (i *Image) ToBoxFormat() string {
	str := fmt.Sprintf("\"%s\": ", i.Name)
	last := len(i.Boxes) - 1
	for ind, box := range i.Boxes {
		str += fmt.Sprintf("(%d, %d, %d, %d)", box[0], box[1], box[2], box[3])
		if ind != last {
			str += ", "
		}
	}

	str += ";"
	return str
}
