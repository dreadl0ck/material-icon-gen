package material_icon_gen

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GenerateIconsPNG(path, url string, hook func(newBase string, color string)) {
	CloneIcons(path, url)

	// rename icons
	_ = os.Mkdir(filepath.Join(path, "renamed"), 0o700)

	files, err := ioutil.ReadDir(filepath.Join(path, "png", "black"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// fmt.Println(f.Name())

		var (
			oldPath = filepath.Join(path, "png", "black", filepath.Base(f.Name()), "twotone-4x.png")
			newBase = filepath.Join(path, "renamed", filepath.Base(f.Name()))
			newPath = newBase + ".png"
		)

		err = os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("renamed", oldPath, "to", newPath)

		GenerateSizes(newBase, newPath)
		if hook != nil {
			hook(newBase, "black")
		}
	}
}

// this will generate a subset of the icons with a different imgType
// call after generateIcons, the image repo needs to be present
func GenerateAdditionalIcons(path string, subset map[string]string, hook func(newBase string, color string)) {
	files, err := ioutil.ReadDir(filepath.Join(path, "png", "black"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// only process files included in the subset
		if imgType, ok := subset[f.Name()]; ok {
			//fmt.Println(f.Name())

			var (
				oldPath = filepath.Join(path, "png", "black", filepath.Base(f.Name()), imgType+"-4x.png")
				newBase = filepath.Join(path, "renamed", filepath.Base(f.Name())+"_"+imgType)
				newPath = newBase + ".png"
			)

			err = os.Rename(oldPath, newPath)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("renamed", oldPath, "to", newPath)

			GenerateSizes(newBase, newPath)
			if hook != nil {
				hook(newBase, "black")
			}
		}
	}
}

func GenerateSizes(newBase string, newPath string) {
	data, err := ioutil.ReadFile(newPath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	for _, size := range []uint{16, 24, 32, 48, 96} {
		newImage := resize.Resize(size, size, img, resize.Lanczos3)

		f, errCreate := os.Create(newBase + strconv.Itoa(int(size)) + ".png")
		if errCreate != nil {
			log.Fatal(errCreate)
		}

		err = png.Encode(f, newImage)
		if err != nil {
			log.Fatal(err)
		}

		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateIcon(text string) {

	var (
		name = strings.ReplaceAll(text, "/", "")
		size = 96.0
	)

	//fmt.Println(name)

	im, err := gg.LoadPNG("/tmp/icons/material-icons-png/renamed/check_box_outline_blank.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(int(size), int(size))
	// dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	var fontSize float64

	switch {
	case len(name) > 12:
		fontSize = 8
	case len(name) > 10:
		fontSize = 11
	case len(name) > 8:
		fontSize = 12
	case len(name) > 6:
		fontSize = 13
	default:
		fontSize = 15
	}

	if err = dc.LoadFontFace(filepath.Join("Roboto", "Roboto-Black.ttf"), fontSize); err != nil {
		panic(err)
	}

	// dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
	dc.DrawImage(im, 0, 0)
	// dc.DrawStringWrapped(text, size/2, size/2, 0.5, 0.5, 80, 1, 0)

	dc.DrawStringAnchored(name, size/2, size/2, 0.5, 0.5)

	dc.Clip()

	var (
		imgBase = filepath.Join("/tmp", "icons", "material-icons-png", "renamed", name)
		imgPath = imgBase + ".png"
	)

	// for testing:
	// imgBase = filepath.Join("/tmp", "icons", "V2", text)
	// imgPath = imgBase + ".png"
	// os.MkdirAll(filepath.Dir(imgBase), 0700)

	err = dc.SavePNG(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	GenerateSizes(imgBase, imgPath)
}