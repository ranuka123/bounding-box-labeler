package image

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	images    []Image
	imagePath string
	boxFile   string = "default.box"
)

func Get(id int) *Image {
	if 0 <= id && id < len(images) {
		return &images[id]
	}
	return nil
}

func Size() int {
	return len(images)
}

func LoadImages(path string) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	images = []Image{}

	for _, file := range files {

		name := file.Name()

		if isImage(name) {
			images = append(images, Image{Name: name})
		}
	}

	if len(images) == 0 {
		return fmt.Errorf("no images found")
	}

	if path[len(path)-1:] != "/" {
		path += "/"
	}
	imagePath = path

	return nil
}

func LoadBoxes(path string) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	targetFileName := ""
	for _, file := range files {

		name := file.Name()

		if strings.Contains(name, ".box") {
			targetFileName = name
			break
		}
	}

	if len(targetFileName) > 0 {

		if err := addBoxesFrom(targetFileName); err != nil {
			return err
		}

		boxFile = targetFileName
	}

	return nil
}

func addBoxesFrom(file string) error {

	boxMapping, err := parseBoxFile(file)
	if err != nil {
		return err
	}

	if len(boxMapping) > 0 {
		for i, img := range images {
			boxes, ok := boxMapping[img.Name]
			if ok {
				images[i].Boxes = boxes
			}
		}
	}

	return nil

}

func isImage(name string) bool {
	return strings.Contains(name, ".jpg") || strings.Contains(name, ".jpeg") || strings.Contains(name, ".png")
}
