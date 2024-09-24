package input

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetFilenamesForPath returns every path+filename found in `path`
func GetFilenamesForPath(path string) ([]string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	images := make([]string, 0)

	fullPath := rootPath + path
	err = filepath.Walk(fullPath, func(imagePath string, info os.FileInfo, err error) error {
		if imagePath == fullPath {
			return nil
		}

		if err != nil {
			return err
		}

		images = append(images, imagePath)

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "could not finish reading files from folder")
	}

	return images, nil
}
