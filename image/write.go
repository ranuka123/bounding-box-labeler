package image

import (
	"bufio"
	"fmt"
	"os"
)

func WriteBoxes() error {
	f, err := os.Create(fmt.Sprintf("%s%s", imagePath, boxFile))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, image := range images {
		if len(image.Boxes) > 0 {
			fmt.Fprintln(w, image.ToBoxFormat())
		}
	}

	return w.Flush()
}
