package output

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type areaProcessorConfig struct {
}

func newAreaProcessor(area *counters.Counter, calculatedAreaHeight int, isLastArea bool) *areaProcessor {
	return &areaProcessor{
		Counter:              area,
		calculatedAreaHeight: calculatedAreaHeight,
		isLastArea:           isLastArea,
	}
}

type areaProcessor struct {
	*counters.Counter
	calculatedAreaHeight int

	areaCanvas *gg.Context

	isLastArea bool
}

func (a *areaProcessor) processArea(card *counters.Card, template *counters.CardsTemplate) (err error) {
	counters.Merge(&a.Settings, card.Settings)

	a.Width = (template.Width) - int(template.Margins*2)
	a.Height = a.calculatedAreaHeight

	if a.areaCanvas, err = counters.GetCanvas(&a.Settings, a.Width, a.Height, template); err != nil {
		return errors.Wrap(err, "error trying to create a canvas")
	}

	if err = a.Images.DrawImagesOnCanvas(&a.Settings, a.areaCanvas, a.Width, a.Height); err != nil {
		return errors.Wrap(err, "error trying to process image")
	}

	if err = a.Texts.DrawTextsOnCanvas(a.Settings, a.areaCanvas, a.Width, a.Height); err != nil {
		return errors.Wrap(err, "error trying to draw text")
	}

	if !a.isLastArea && a.Frame {
		a.drawFrame(a.BorderWidth, a.BorderColor)
	}

	return nil
}

func (a *areaProcessor) drawFrame(w float64, col color.Color) {
	a.areaCanvas.Push()
	a.areaCanvas.SetColor(col)
	a.areaCanvas.SetLineWidth(w)
	frameX := float64(a.areaCanvas.Width())
	frameY := float64(a.areaCanvas.Height())
	a.areaCanvas.DrawRectangle(0, 0, frameX, frameY)
	a.areaCanvas.Stroke()
	a.areaCanvas.Pop()
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
