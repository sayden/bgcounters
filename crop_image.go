package counters

import (
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

// CropFolderToContent is like CropToContentFile but it will crop
// all images in the provided folder
func CropFolderToContent(folderPath string) error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := rootPath + string(os.PathSeparator) + folderPath
	return filepath.Walk(fullPath, func(imagePath string, info os.FileInfo, err error) error {
		if imagePath == fullPath {
			return nil
		}

		if err != nil {
			return err
		}

		err = CropToContentFile(imagePath)
		if err != nil {
			log.WithError(err).Error("omitting file")
		}

		return nil
	})
}

// CropToContentFile is like CropToContent but the file read will be
// overriden with the resulting image. USE WITH CARE!
func CropToContentFile(imagePath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return errors.Wrapf(err, "error trying to open file '%s'", imagePath)
	}

	newImage := CropToContent(img)
	if newImage == img {
		return nil
	}

	return imaging.Save(newImage, imagePath)
}

func CropToContent(i image.Image) image.Image {
	bounds := i.Bounds()
	max := bounds.Max
	min := bounds.Min

	topLimit := min.Y
	leftLimit := max.X
	lowerLimit := 0
	rightLimit := 0

	var alpha uint32

topLimitLoop:
	for row := 0; row < max.Y; row++ {
		for column := 0; column < max.X; column++ {
			_, _, _, alpha = i.At(column, row).RGBA()
			//If color is found
			if alpha > 0 {
				if topLimit == 0 {
					if topLimit > 0 {
						topLimit = row - 1
					} else {
						topLimit = row
					}
				}
				break topLimitLoop
			}
		}
	}

leftLimitLoop:
	for column := 0; column < max.X; column++ {
		for row := topLimit; row < max.Y; row++ {
			_, _, _, alpha = i.At(column, row).RGBA()
			//If color is found
			if alpha > 0 {
				leftLimit = column
				break leftLimitLoop
			}
		}
	}
rightLimitLoop:
	for column := max.X; column >= leftLimit; column-- {
		for row := topLimit; row < max.Y; row++ {
			_, _, _, alpha = i.At(column, row).RGBA()
			//If color is found
			if alpha > 0 {
				rightLimit = column + 1
				break rightLimitLoop
			}
		}
	}

lowerLimitLoop:
	for row := max.Y; row >= topLimit; row-- {
		for column := leftLimit; column < rightLimit; column++ {
			_, _, _, alpha = i.At(column, row).RGBA()
			//If color is found
			if alpha > 0 {
				lowerLimit = row + 1
				break lowerLimitLoop
			}
		}
	}

	rect := image.Rect(leftLimit, topLimit, rightLimit, lowerLimit)
	if rect == bounds {
		return i
	}

	newImg := imaging.Crop(i, rect)

	return newImg
}
