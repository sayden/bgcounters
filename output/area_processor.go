package output

import (
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type areaProcessorConfig struct {
	area                 *counters.Counter
	calculatedAreaHeight int
	isLastArea           bool
}

func newAreaProcessor(c *areaProcessorConfig) *areaProcessor {
	return &areaProcessor{
		Counter:              c.area,
		calculatedAreaHeight: c.calculatedAreaHeight,
		isLastArea:           c.isLastArea,
	}
}

type areaProcessor struct {
	*counters.Counter
	calculatedAreaHeight int

	areaCanvas *gg.Context

	isLastArea bool
}

// processArea draw images and texts into a new canvas, stored in a.areaCanvas
func (a *areaProcessor) processArea(card *counters.Card, template *counters.CardsTemplate) error {
	counters.Merge(&a.Settings, card.Settings)

	a.Width = (template.Width) - int(template.Margins*2)
	a.Height = a.calculatedAreaHeight

	var err error
	if a.areaCanvas, err = counters.GetCanvas(&a.Settings, a.Width, a.Height, template); err != nil {
		return errors.Wrap(err, "error trying to create a canvas")
	}

	if err = counters.DrawImagesOnCanvas(a.Images, &a.Settings, a.areaCanvas, a.Width, a.Height); err != nil {
		return errors.Wrap(err, "error trying to process image")
	}

	counters.DrawTextsOnCanvas(a.Texts, a.Settings, a.areaCanvas, a.Width, a.Height)

	if !a.isLastArea && a.Frame {
		drawFrame(a.areaCanvas, a.BorderWidth, a.BorderColor)
	}

	if !a.Settings.SkipBorders {
		a.drawBorders()
	}

	return nil
}

// drawOnCard draw the area canvas on the card canvas
func (a *areaProcessor) drawOnCard(template *counters.CardsTemplate, cardCanvas *gg.Context, x, y float64) error {
	cardCanvas.DrawImage(a.areaCanvas.Image(), int(x), int(y))

	if template.DrawGuides {
		guidesImage, err := counters.DrawGuides(a.Settings)
		if err != nil {
			return errors.Wrap(err, "error tyring to draw guides")
		}
		cardCanvas.DrawImage(*guidesImage, int(x), int(y))
	}

	return nil
}

func (a *areaProcessor) drawBorders() {
	a.areaCanvas.Push()
	a.areaCanvas.SetColor(a.Settings.BorderColor)
	a.areaCanvas.SetLineWidth(a.Settings.BorderWidth)
	a.areaCanvas.DrawRectangle(0, 0, float64(a.Settings.Width), float64(a.Settings.Height))
	a.areaCanvas.Stroke()
	a.areaCanvas.Pop()
}
