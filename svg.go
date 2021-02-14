package material_icon_gen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	svgcheck "github.com/h2non/go-is-svg"
)

// MaterialIconSVG models the material icon svg structure
type MaterialIconSVG struct {
	XMLName xml.Name `xml:"svg"`

	Xmlns   string `xml:"xmlns,attr"`
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
	Paths   []Path `xml:"path"`

	Rect *struct {
		Text        string `xml:",chardata"`
		X           string `xml:"x,attr"`
		Y           string `xml:"y,attr"`
		Width       string `xml:"width,attr"`
		Height      string `xml:"height,attr"`
		Stroke      string `xml:"stroke,attr"`
		StrokeWidth string `xml:"stroke-width,attr"`
		Fill        string `xml:"fill,attr"`
	} `xml:"rect,omitempty"`

	Text *struct {
		Text             string `xml:",chardata"`
		X                string `xml:"x,attr"`
		Y                string `xml:"y,attr"`
		DominantBaseline string `xml:"dominant-baseline,attr"`
		TextAnchor       string `xml:"text-anchor,attr"`
	} `xml:"text,omitempty"`
}

type Path struct {
	Text    string `xml:",chardata"`
	Opacity string `xml:"opacity,attr"`
	D       string `xml:"d,attr"`
	Style   string `xml:"style,attr,omitempty"`
}

func (s *MaterialIconSVG) ResizeSVG(width, height int) {
	s.Height = strconv.Itoa(height)
	s.Width = strconv.Itoa(width)
}


func GenerateIconsSVG(path string, url string, sizes []int, coloredIcons map[string][]string, hook func(newBase string, color string)) {
	CloneIcons(path, url)

	// rename icons
	_ = os.Mkdir(filepath.Join(path, "renamed"), 0o700)

	files, err := ioutil.ReadDir(filepath.Join(path, "svg"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// fmt.Println(f.Name())

		var (
			oldPath = filepath.Join(path, "svg", filepath.Base(f.Name()), "twotone.svg")
			newBase = filepath.Join(path, "renamed", filepath.Base(f.Name()))
			newPath = newBase + ".svg"
		)

		err = copyFile(oldPath, newPath)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("renamed", oldPath, "to", newPath)

		if colorNames, ok := coloredIcons[f.Name()]; ok {
			for _, c := range colorNames {
				GenerateSizesSVG(newBase, newPath, c, sizes)
				if hook != nil {
					hook(newBase, c)
				}
			}
		} else {
			GenerateSizesSVG(newBase, newPath, "black", sizes)
			if hook != nil {
				hook(newBase, "black")
			}
		}
	}
}

func GenerateAdditionalIconsSVG(path string, sizes []int, subset map[string]string, hook func(newBase string, color string)) {
	files, err := ioutil.ReadDir(filepath.Join(path, "svg"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// only process files included in the subset
		if imgType, ok := subset[f.Name()]; ok {
			// fmt.Println(f.Name())

			var (
				oldPath = filepath.Join(path, "svg", filepath.Base(f.Name()), imgType+".svg")
				newBase = filepath.Join(path, "renamed", filepath.Base(f.Name())+"_"+imgType)
				newPath = newBase + ".svg"
			)

			err = copyFile(oldPath, newPath)
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Println("renamed", oldPath, "to", newPath)

			GenerateSizesSVG(newBase, newPath, "black", sizes)
			if hook != nil {
				hook(newBase, "black")
			}
		}
	}
}

func GenerateSizesSVG(newBase string, newPath string, color string, sizes []int) {
	svgFile, err := os.Open(newPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer func() {
		errClose := svgFile.Close()
		if errClose != nil {
			fmt.Println(errClose)
		}
	}()

	s := new(MaterialIconSVG)
	if err = xml.NewDecoder(svgFile).Decode(&s); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to parse (%v)\n", err)
		return
	}

	for i := range s.Paths {
		s.Paths[i].Style = "fill: " + color + ";"
	}

	if s.ViewBox == "" {
		s.ViewBox = "0 0 100 100"
	}

	for _, size := range sizes {

		s.ResizeSVG(size, size)
		f, errCreate := os.Create(newBase + "_" + color + strconv.Itoa(size) + ".svg")
		if errCreate != nil {
			log.Fatal(errCreate)
		}

		var buf bytes.Buffer

		if err = xml.NewEncoder(io.MultiWriter(f, &buf)).Encode(s); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to encode (%v)\n", err)
			return
		}

		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}

		if !svgcheck.Is(buf.Bytes()) {
			log.Fatal("invalid SVG image generated:", f.Name())
		}
	}
}

func GenerateIconSVG(path string, text string, sizes []int) string {

	var (
		name = strings.ReplaceAll(text, "/", "")
		size = 96
	)

	//fmt.Println(name)

	var (
		x = `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="` + strconv.Itoa(size) + `" height="` + strconv.Itoa(size) + `">
<rect x="0" y="0" width="` + strconv.Itoa(size) + `" height="` + strconv.Itoa(size) + `" stroke="red" stroke-width="3px" fill="white"/>
<text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle">` + name + `</text>
</svg>`
	)

	var (
		imgBase = filepath.Join(path, "renamed", name)
		imgPath = imgBase + ".svg"
	)

	file, err := os.Create(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errClose := file.Close()
		if errClose != nil {
			fmt.Println(errClose)
		}
	}()

	_, err = file.WriteString(x)
	if err != nil {
		log.Fatal(err)
	}

	if !svgcheck.Is([]byte(x)) {
		log.Fatal("invalid SVG", x)
	}

	GenerateSizesSVG(imgBase, imgPath, "black", sizes)

	return imgBase
}
