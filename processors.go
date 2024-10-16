package counters

import (
	"github.com/fogleman/gg"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func Merge(d *Settings, s2 Settings, opt ...func(*mergo.Config)) {
	d.FontColor = nil
	d.StrokeColor = nil
	d.BorderColor = nil
	d.BgColor = nil

	mergo.Merge(d, s2, opt...)

	SetColors(d)
}

// DrawImagesOnCanvas using the provided height `h` and width `w`
func DrawImagesOnCanvas(images []Image, s *Settings, areaCanvas *gg.Context, w, h int) error {
	// sort the internal objects by position
	slices.SortFunc(images, func(i, j Image) int {
		return i.Position - j.Position
	})

	for _, image := range images {
		Merge(&image.Settings, *s)

		image.Width = w
		image.Height = h

		if err := image.Draw(areaCanvas, image.Position, image.Settings); err != nil {
			return errors.Wrap(err, "could not draw image")
		}
	}

	return nil
}

// DrawTextsOnCanvas draws the texts provided on areaCanvas at positions `w` and `h`
// using the provided Settings
func DrawTextsOnCanvas(texts []Text, s Settings, areaCanvas *gg.Context, w, h int) {
	for _, text := range texts {
		Merge(&text.Settings, s)

		text.Width = w
		text.Height = h
		text.Draw(areaCanvas, text.Position, text.Settings)
	}
}
