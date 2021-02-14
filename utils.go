package material_icon_gen

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
)

const (
	DefaultSvgURL = "https://github.com/dreadl0ck/material-icons.git"
	DefaultPngURL = "https://github.com/dreadl0ck/material-icons-png.git"
)

func CloneIcons(path string, url string) {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		fmt.Println("repo exists, pulling")
		r, err := git.PlainOpen(path)
		if err != nil {
			log.Fatal(err)
		}

		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		err = w.Pull(&git.PullOptions{})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate){
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("cloned icon repository to", path)
	}
}
