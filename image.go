package counters

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/thehivecorporation/log"
)

type Image struct {
	Settings
	Positioner

	Path          string  `json:"path"`
	Scale         float64 `json:"scale,omitempty" default:"1"`
	AvoidCropping bool    `json:"avoid_cropping,omitempty"`
}

func (i *Image) SetWidth(w int) {
	i.Width = w
}

func (i *Image) SetHeight(h int) {
	i.Height = h
}

func (i *Image) GetSettings() Settings {
	return i.Settings
}

func (i *Image) Draw(dc *gg.Context, pos int, _ Settings) error {
	img, err := gg.LoadImage(i.Path)
	if err != nil {
		log.WithField("image", i.Path).Error("error trying to load image in 'Image' item")
		return err
	}

	if !i.AvoidCropping {
		img = CropToContent(img)
	}

	if i.Rotation != 0 {
		img = imaging.Rotate(img, i.Rotation, color.Transparent)
	}

	if i.ImageScaling == "" {
		i.ImageScaling = IMAGE_SCALING_FIT_NONE
	}

	if i.Scale == 0 {
		i.Scale = 1
	}

	switch i.ImageScaling {
	case IMAGE_SCALING_FIT_WIDTH:
		img = imaging.Resize(img, int(math.Ceil(float64(dc.Width())*i.Scale)), 0, imaging.Box)
	case IMAGE_SCALING_FIT_WRAP:
		img = imaging.Resize(img, dc.Width(), dc.Height(), imaging.Box)
	case IMAGE_SCALING_FIT_HEIGHT:
		img = imaging.Resize(img, int(math.Ceil(float64(dc.Height())*i.Scale)), 0, imaging.Box)
	case IMAGE_SCALING_FIT_NONE:
		// Do nothing, image untouched from original
	default:
		// Do nothing, image untouched from original
	}

	x, y, ax, ay, err := i.getObjectPositions(pos, i.Settings)
	if err != nil {
		return err
	}
	if i.ShadowDistance != 0 {
		shadow := getShadowFromImage(img, i.ShadowDistance, SIGMA)
		x1 := math.Floor(x + float64(i.ShadowDistance))
		y1 := math.Floor(y + float64(i.ShadowDistance))
		dc.DrawImageAnchored(shadow, int(x1), int(y1), ax, ay)
	}

	dc.DrawImageAnchored(img, int(x), int(y), ax, ay)

	return nil
}

type ImagePrototype struct {
	Image

	PathList []string `json:"path_list"`
}

func getShadowFromImage(img image.Image, shadowDistance int, sigma int) image.Image {
	grey := imaging.AdjustBrightness(img, -100)
	w, h := getShadowImageSize(img, shadowDistance, sigma)
	temp := gg.NewContext(w*2, h*2)
	temp.DrawImageAnchored(grey, w/2, h/2, 0.5, 0.5)
	grey = imaging.AdjustBrightness(temp.Image(), 15)
	grey = imaging.Blur(grey, float64(sigma))

	return CropToContent(grey)
}

func getShadowImageSize(img image.Image, shadowDistance int, sigma int) (int, int) {
	rect := img.Bounds()
	w := rect.Dx() + sigma*sigma
	h := rect.Dy() + sigma*sigma

	return w, h
}

func applyImageScaling(i *Image, scaling float64) {
	i.Margins *= scaling

	i.FontHeight *= scaling

	i.ShadowDistance = int(scaling * float64(i.ShadowDistance))

	i.BorderWidth *= scaling
	if i.BorderWidth < 1 {
		i.BorderWidth = 1
	}

	i.XShift *= scaling
	i.YShift *= scaling

	i.StrokeWidth *= scaling

	// if i.Scale == 0 {
	// 	i.Scale = scaling
	// } else {
	// 	i.Scale *= scaling
	// }
}
