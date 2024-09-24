package output

import (
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

func CountersToBlocks(frontTemplate, backTemplate *counters.CounterTemplate) error {
	var canvas *gg.Context
	if frontTemplate.Mode == counters.TEMPLATE_MODE_TILES {
		return errors.New("tiles cannot be set when producing blocks")
	}

	counterPos := 0
	row := 0
	fileNumber := 1
	//Iterate rows and columns, painting a counter on each
iteration:
	for {
		for columns := 0; columns < frontTemplate.Columns; columns++ {
			if counterPos == len(frontTemplate.Counters) {
				break iteration
			}

			counter := frontTemplate.Counters[counterPos]

			counterCanvas, err := getCounterCanvas(&counter, frontTemplate)
			if err != nil {
				return err
			}

			blockCanvas, err := getBlockCanvasFromCounterCanvas(counterCanvas, &counter)
			if err != nil {
				return err
			}

			if backTemplate != nil {
				backCounter := backTemplate.Counters[counterPos]
				backCounterCanvas, err := getCounterCanvas(&backCounter, backTemplate)
				if err != nil {
					return err
				}

				addBackCounterToBlockCanvas(backCounterCanvas, blockCanvas)
			}

			if err = writeCounterToFile(blockCanvas, counter, frontTemplate, &fileNumber, &columns, &row, canvas); err != nil {
				return err
			}

			counterPos++
		}

		row++
	}

	// Save result on a file
	if frontTemplate.Mode == counters.TEMPLATE_MODE_TILES {
		return canvas.SavePNG(frontTemplate.OutputFolder)
	}

	return nil
}

func addBackCounterToBlockCanvas(backCounter, blockCanvas *gg.Context) {
	canvasWidth := blockCanvas.Width()
	canvasHeight := blockCanvas.Height()
	margin := int(float64(canvasWidth) * 0.0508)

	img := backCounter.Image()
	img = imaging.Rotate(img, 180, color.Black)

	blockCanvas.DrawImage(img, int(canvasWidth-margin-backCounter.Width()), int(canvasHeight-margin-backCounter.Height()))
}

func getBlockCanvasFromCounterCanvas(counterCanvas *gg.Context, cc *counters.Counter) (*gg.Context, error) {
	const canvasRatio = 2.49
	canvasWidth := float64(counterCanvas.Width()) * canvasRatio
	canvasHeigth := float64(counterCanvas.Height()) * canvasRatio
	dc := gg.NewContext(int(canvasWidth), int(canvasHeigth))
	if err := dc.LoadFontFace(cc.FontPath, cc.FontHeight); err != nil {
		return nil, err
	}

	dc.Push()
	dc.SetColor(cc.BgColor)
	dc.DrawRectangle(0, 0, canvasWidth, canvasHeigth)
	dc.Fill()
	dc.Pop()

	if cc.FontColorS != "" {
		// counters.GetValidColorForString(cc.FontColorS, t.BgColor)
		counters.GetValidColorForString(cc.FontColorS, cc.BgColor)
	}

	// Draw the counter into the new canvas
	margin := canvasWidth * 0.0508
	bottomMargin := int(canvasHeigth - margin)
	dc.DrawImage(counterCanvas.Image(), int(margin), bottomMargin-cc.Height)

	return dc, nil
}

func getCounterCanvas(counter *counters.Counter, template *counters.CounterTemplate) (*gg.Context, error) {
	dc := GetCanvasForCounter(counter, template)

	// Draw background image
	err := drawBackgroundImage(dc, counter.Settings)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to draw background image")
	}

	// Draw images
	if err = counters.DrawImagesOnCanvas(counter.Images, &counter.Settings, dc, counter.Width, counter.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to process image")
	}

	// Draw texts
	counters.DrawTextsOnCanvas(counter.Texts, counter.Settings, dc, counter.Width, counter.Height)

	// Draw guides
	if template.DrawGuides {
		guides, err := counters.DrawGuides(counter.Settings)
		if err != nil {
			return nil, err
		}
		dc.DrawImage(*guides, 0, 0)
	}

	return dc, nil
}
