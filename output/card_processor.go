package output

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/sayden/counters"
)

type cardProcessorConfig struct {
	template   *counters.CardsTemplate
	cardCanvas *gg.Context
}

func newCardProcessor(cfg *cardProcessorConfig) *cardProcessor {
	return &cardProcessor{
		template:   cfg.template,
		cardCanvas: cfg.cardCanvas,
	}
}

type cardProcessor struct {
	template *counters.CardsTemplate

	cardCanvas *gg.Context
}

/*
processCard processes a single card by merging its settings with the template settings,
drawing the background image, processing each area of the card, drawing texts, and optionally
drawing borders. Finally, it draws the processed card onto the provided sheet at the specified
column and row position.

Parameters:
- sheet: The context of the sheet where the card will be drawn.
- card: The card to be processed.
- columns: The column position on the sheet where the card will be drawn.
- rows: The row position on the sheet where the card will be drawn.
*/
func (c *cardProcessor) processCard(sheet *gg.Context, card *counters.Card, columns, rows int) error {
	counters.Merge(&card.Settings, c.template.Settings)
	counters.SetColors(&card.Settings)

	var err error
	c.cardCanvas, err = counters.GetCanvas(&card.Settings, c.template.Width, c.template.Height, c.template)
	if err != nil {
		return err
	}

	if err = card.DrawBackgroundImage(c.cardCanvas); err != nil {
		return err
	}

	err = card.Images.DrawImagesOnCanvas(&card.Settings, c.cardCanvas, card.Width, card.Height)
	if err != nil {
		return err
	}
	// Height when all areas have the same height
	numberOfAreas := len(card.Areas)

	//calculatedAreaHeight := (float64(template.Height) - (template.Margins * 2)) / float64(numberOfAreas)

	areasHeights := getAreasHeights(card.Areas, card.Height, card.Margins)
	// Process each area on the text
	y := c.template.Margins

	for areaIndex, area := range card.Areas {
		isLastAreaOfCard := areaIndex != numberOfAreas
		card.Areas[areaIndex].Height = int(math.Floor(areasHeights[areaIndex]))

		areaProc := newAreaProcessor(&areaProcessorConfig{
			area:                 &area,
			calculatedAreaHeight: int(math.Floor(areasHeights[areaIndex])),
			isLastArea:           isLastAreaOfCard,
		})
		if err = areaProc.processArea(card, c.template); err != nil {
			return err
		}

		x := c.template.Margins
		if err = areaProc.drawOnCard(c.template, c.cardCanvas, x, y); err != nil {
			return err
		}
		y += areasHeights[areaIndex]
	}

	if err = card.Texts.DrawTextsOnCanvas(card.Settings, c.cardCanvas, card.Width, card.Height); err != nil {
		return err
	}

	c.maybeDrawBorders(card)

	sheet.DrawImage(c.cardCanvas.Image(), columns*card.Width, rows*card.Height)

	return nil
}

func getAreasHeights(areas []counters.Counter, parentHeight int, margins float64) (hs []float64) {
	hs = make([]float64, len(areas))
	availableH := float64(parentHeight) - (margins * 2)
	hasCustomHeight := make([]bool, len(areas))
	totalNonCustomAreas := 0
	for i, area := range areas {
		if area.Height != parentHeight {
			hasCustomHeight[i] = true
			availableH -= float64(area.Height)
			continue
		}
		totalNonCustomAreas++
	}

	availableSpaceForNonCustom := availableH / float64(totalNonCustomAreas)
	for i, isCustom := range hasCustomHeight {
		if isCustom {
			hs[i] = float64(areas[i].Height)
		} else {
			hs[i] = availableSpaceForNonCustom
		}
	}

	return
}

func (c *cardProcessor) maybeDrawBorders(card *counters.Card) {
	borderColorIsSet := card.Settings.BorderColor != nil
	borderWidthIsSet := card.Settings.BorderWidth != 0

	if borderColorIsSet && borderWidthIsSet {
		c.cardCanvas.Push()
		c.cardCanvas.SetColor(card.Settings.BorderColor)
		c.cardCanvas.SetLineWidth(card.Settings.BorderWidth)
		c.cardCanvas.DrawRectangle(0, 0, float64(card.Settings.Width), float64(card.Settings.Height))
		c.cardCanvas.Stroke()
		c.cardCanvas.Pop()
	}
}
