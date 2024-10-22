package counters

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type Card struct {
	Settings
	Areas  []Counter `json:"areas,omitempty"`
	Texts  Texts     `json:"texts,omitempty"`
	Images Images    `json:"images,omitempty"`
}

func (c *Card) Image(template *CardsTemplate) (image.Image, error) {
	cardCanvas, err := c.ToCanvas(template)
	if err != nil {
		return nil, err
	}

	return cardCanvas.Image(), nil
}

func (c *Card) ToCanvas(template *CardsTemplate) (*gg.Context, error) {
	Merge(&c.Settings, template.Settings)
	SetColors(&c.Settings)

	cardCanvas, err := GetCanvas(&c.Settings, template.Width, template.Height, template)
	if err != nil {
		return nil, err
	}

	if err = c.DrawBackgroundImage(cardCanvas); err != nil {
		return nil, err
	}

	if err = c.Images.DrawImagesOnCanvas(&c.Settings, cardCanvas, c.Width, c.Height); err != nil {
		return nil, err
	}
	// Height when all areas have the same height
	numberOfAreas := len(c.Areas)

	//calculatedAreaHeight := (float64(template.Height) - (template.Margins * 2)) / float64(numberOfAreas)

	areasHeights := c.GetAreasHeights()
	// Process each area on the text
	y := template.Margins

	for areaIndex, area := range c.Areas {
		isLastAreaOfCard := areaIndex != numberOfAreas
		c.Areas[areaIndex].Height = int(math.Floor(areasHeights[areaIndex]))

		// area.Width = (template.Width) - int(template.Margins*2)
		// areaCanvas, err := c.processAreav2(template, &area, c.Areas[areaIndex].Height, isLastAreaOfCard)
		areaCanvas, err := c.processAreav2(&area, template, int(math.Floor(areasHeights[areaIndex])), isLastAreaOfCard)
		// areaCanvas, err := area.Canvas(false)
		if err != nil {
			return nil, err
		}

		// TODO: Remove this
		f, _ := os.Create(fmt.Sprintf("/tmp/area_%d.png", areaIndex))
		defer f.Close()
		areaCanvas.EncodePNG(f)

		x := template.Margins
		if err = c.drawOnCard(template, cardCanvas, areaCanvas, x, y); err != nil {
			return nil, err
		}
		y += areasHeights[areaIndex]
	}

	if err = c.Texts.DrawTextsOnCanvas(c.Settings, cardCanvas, c.Width, c.Height); err != nil {
		return nil, err
	}

	c.maybeDrawBorders(cardCanvas)

	return cardCanvas, nil
}

func (c *Card) EncodeImage(w io.Writer, t *CardsTemplate) error {
	cardCanvas, err := c.ToCanvas(t)
	if err != nil {
		return err
	}

	if err = cardCanvas.EncodePNG(w); err != nil {
		return err
	}

	return nil
}

func (c *Card) processAreav2(area *Counter, t *CardsTemplate, calculatedAreaHeight int, isLastArea bool) (*gg.Context, error) {
	Merge(&area.Settings, c.Settings)

	area.Width = (t.Width) - int(t.Margins*2)
	area.Height = calculatedAreaHeight

	areaCanvas, err := area.Canvas(false)
	if err != nil {
		return nil, err
	}

	if !isLastArea && area.Frame {
		c.drawFrame(areaCanvas, c.BorderWidth, c.BorderColor)
	}

	return areaCanvas, nil
}

func (c *Card) processArea(t *CardsTemplate, area *Counter, calculatedAreaHeight int, isLastArea bool) (*gg.Context, error) {
	Merge(&area.Settings, c.Settings)

	c.Width = (t.Width) - int(t.Margins*2)
	c.Height = calculatedAreaHeight

	areaCanvas, err := GetCanvas(&c.Settings, c.Width, c.Height, t)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to create a canvas")
	}

	if err = c.Images.DrawImagesOnCanvas(&c.Settings, areaCanvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to process image")
	}

	if err = c.Texts.DrawTextsOnCanvas(c.Settings, areaCanvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to draw text")
	}

	if !isLastArea && area.Frame {
		c.drawFrame(areaCanvas, c.BorderWidth, c.BorderColor)
	}

	return areaCanvas, nil
}

func (a *Card) drawOnCard(t *CardsTemplate, cardCanvas, areaCanvas *gg.Context, x, y float64) error {
	cardCanvas.DrawImage(areaCanvas.Image(), int(x), int(y))

	if t.DrawGuides {
		guidesImage, err := DrawGuides(a.Settings)
		if err != nil {
			return errors.Wrap(err, "error tyring to draw guides")
		}
		cardCanvas.DrawImage(*guidesImage, int(x), int(y))
	}

	return nil
}

func (a *Card) drawFrame(areaCanvas *gg.Context, w float64, col color.Color) {
	areaCanvas.Push()
	areaCanvas.SetColor(col)
	areaCanvas.SetLineWidth(w)
	frameX := float64(areaCanvas.Width())
	frameY := float64(areaCanvas.Height())
	areaCanvas.DrawRectangle(0, 0, frameX, frameY)
	areaCanvas.Stroke()
	areaCanvas.Pop()
}

func (c *Card) maybeDrawBorders(cardCanvas *gg.Context) {
	borderColorIsSet := c.Settings.BorderColor != nil
	borderWidthIsSet := c.Settings.BorderWidth != 0

	if borderColorIsSet && borderWidthIsSet {
		cardCanvas.Push()
		cardCanvas.SetColor(c.Settings.BorderColor)
		cardCanvas.SetLineWidth(c.Settings.BorderWidth)
		cardCanvas.DrawRectangle(0, 0, float64(c.Settings.Width), float64(c.Settings.Height))
		cardCanvas.Stroke()
		cardCanvas.Pop()
	}
}

func (c *Card) GetAreasHeights() (hs []float64) {
	hs = make([]float64, len(c.Areas))
	availableH := float64(c.Height) - (c.Margins * 2)
	hasCustomHeight := make([]bool, len(c.Areas))
	totalNonCustomAreas := 0
	for i, area := range c.Areas {
		if area.Height == 0 {
			totalNonCustomAreas++
			continue
		}
		if area.Height != c.Height {
			hasCustomHeight[i] = true
			availableH -= float64(area.Height)
			continue
		}
		totalNonCustomAreas++
	}

	availableSpaceForNonCustom := availableH / float64(totalNonCustomAreas)
	for i, isCustom := range hasCustomHeight {
		if isCustom {
			hs[i] = float64(c.Areas[i].Height)
		} else {
			hs[i] = availableSpaceForNonCustom
		}
	}

	return
}

// GetCanvas returns a Canvas with attributes (like background color or size)
// taken from `settings`
func GetCanvas(settings *Settings, width, height int, t *CardsTemplate) (*gg.Context, error) {
	dc := gg.NewContext(width, height)
	err := dc.LoadFontFace(settings.FontPath, settings.FontHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not load font face")
	}

	if settings.BgColor != nil {
		dc.Push()
		dc.SetColor(settings.BgColor)
		dc.DrawRectangle(0, 0, float64(settings.Width), float64(settings.Height))
		dc.Fill()
		dc.Pop()
	}

	if settings.FontColorS != "" {
		ColorFromStringOrDefault(settings.FontColorS, t.BgColor)
	}

	return dc, nil
}

func ApplyCardScaling(t *CardsTemplate) {
	for i := range t.Cards {
		c := t.Cards[i]
		c.Settings.ApplySettingsScaling(t.Scaling)

		for j := range c.Areas {
			ApplyCounterScaling(&c.Areas[j], t.Scaling)
		}

		for j := range c.Images {
			applyImageScaling(&c.Images[j], t.Scaling)
		}

		for j := range c.Texts {
			c.Texts[j].Settings.ApplySettingsScaling(t.Scaling)
		}
	}

	t.Settings.ApplySettingsScaling(t.Scaling)
}
