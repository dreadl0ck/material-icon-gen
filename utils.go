package material_icon_gen

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"os"
)

const (
	svgIconPath = "/tmp/icons/material-icons"
	pngIconPath = "/tmp/icons/material-icons-png"
)

func CloneIcons() {
	_ = os.RemoveAll(pngIconPath)

	_, err := git.PlainClone(pngIconPath, false, &git.CloneOptions{
		URL:      "https://github.com/dreadl0ck/material-icons-png.git",
		Progress: os.Stdout,
	})

	if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
		log.Fatal(err)
	}

	fmt.Println("cloned icon repository to", pngIconPath)
}
