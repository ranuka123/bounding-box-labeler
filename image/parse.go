package image

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var boxInBracket *regexp.Regexp = regexp.MustCompile(`\((.*?)\)`)

func parseBoxFile(fileName string) (map[string][]Box, error) {

	f, err := os.Open(fmt.Sprintf("%s%s", imagePath, fileName))
	if err != nil {
		return map[string][]Box{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	imagesMap := map[string][]Box{}

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), ":")
		fileName = trimDoubleQuotes(split[0])
		boxes := parseBoxes(split[1])

		imagesMap[fileName] = boxes
	}

	if err := scanner.Err(); err != nil {
		return map[string][]Box{}, err
	}

	return imagesMap, nil
}

func trimDoubleQuotes(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}

func parseBoxes(input string) []Box {

	extracted := boxInBracket.FindAll([]byte(input), -1)

	r := strings.NewReplacer("(", "", ")", "", " ", "")

	boxes := []Box{}
	for _, byteStr := range extracted {

		str := r.Replace(string(byteStr))
		boxPoints := strings.Split(str, ",")

		box := Box{}
		for i := 0; i < 4; i++ {
			p, err := strconv.Atoi(boxPoints[i])
			if err != nil {
				log.Println("bad str", input, err)
				return []Box{}
			}
			box[i] = p
		}
		boxes = append(boxes, box)
	}

	return boxes

}
