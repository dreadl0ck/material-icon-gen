package material_icon_gen

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
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
		fmt.Println("material icon repository exists, pulling")
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
		fmt.Println("cloned material icon repository to", path)
	}
}

// copyFile the source file contents to destination
// file attributes wont be copied and an existing file will be overwritten.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if errClose := in.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}