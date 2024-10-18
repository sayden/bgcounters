package counters

import (
	"image"

	"github.com/disintegration/imaging"
)

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
