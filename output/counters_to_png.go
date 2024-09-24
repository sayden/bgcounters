package output

import (
	"os"
	"path"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

func CountersToPNG(template *counters.CounterTemplate) error {
	var canvas *gg.Context
	if template.Mode == counters.TEMPLATE_MODE_TILES {
		// Create canvas
		canvas = gg.NewContext(template.Columns*template.Width, template.Rows*template.Height)
		if err := canvas.LoadFontFace(counters.DEFAULT_FONT_PATH, template.FontHeight); err != nil {
			log.WithField("font_path", counters.DEFAULT_FONT_PATH).Error(err)
			return errors.Wrap(err, "could not load font")
		}
	}

	counterPos := 0
	row := 0
	fileNumber := 1
	//Iterate rows and columns, painting a counter on each
iteration:
	for {
		for columns := 0; columns < template.Columns; columns++ {
			if counterPos == len(template.Counters) {
				break iteration
			}

			counter := template.Counters[counterPos]
			if counter.Skip {
				counterPos++
				continue
			}

			counterCanvas, err := getCounterCanvas(&counter, template)
			if err != nil {
				return err
			}

			if err = writeCounterToFile(counterCanvas, counter, template, &fileNumber, &columns, &row, canvas); err != nil {
				return err
			}

			counterPos++
		}

		row++
	}

	// Save result on a file
	if template.Mode == counters.TEMPLATE_MODE_TILES {
		return canvas.SavePNG(template.OutputFolder)
	}

	return nil
}

type imageWriter struct {
	dc       *gg.Context
	template *counters.CounterTemplate
}

func writeCounterToFile(dc *gg.Context, counter counters.Counter, template *counters.CounterTemplate, fileNumber *int, columns, row *int, canvas *gg.Context) error {
	iw := imageWriter{
		dc:       dc,
		template: template,
	}
	var suffix string
	if counter.Back != nil {
		suffix = "-front"
	}

	switch template.Mode {
	case counters.TEMPLATE_MODE_TEMPLATE:
		n, err := iw.createFile(&counter, fileNumber, suffix)
		if err != nil {
			return errors.Wrap(err, "error trying to write counter file")
		}

		if counter.Back != nil {
			suffix = "-back"
			*fileNumber -= n
			t, err := iw.createFile(counter.Back, fileNumber, suffix)
			if err != nil || t != n {
				return errors.Wrap(err, "error trying to write counter file")
			}
		}
	case counters.TEMPLATE_MODE_TILES:
		if counter.Back != nil {
			iw.createTiledFile(counter.Back, columns, row, canvas)
		}
		iw.createTiledFile(&counter, columns, row, canvas)
	default:
		return errors.New("unknown template mode or template mode is empty")
	}

	return nil
}

// createFile creates a file with the counter image. Filenumber is the filename, a pointer is passed to be able to use
// the multiplier to create more than one file with the same counter
func (iw *imageWriter) createFile(counter *counters.Counter, fileNumber *int, suffix string) (int, error) {
	_ = os.MkdirAll(iw.template.OutputFolder, 0750)

	if counter.Multiplier == 0 {
		counter.Multiplier = 1
	}

	starting := *fileNumber

	// Use sequencing of numbers or a position in the counter texts to name files
	for i := 0; i < counter.Multiplier; i++ {
		filepath := path.Join(iw.template.OutputFolder, counter.GetCounterFilename(iw.template.IndexNumberForFilename, *fileNumber, suffix))
		if counter.Skip {
			continue
		}
		if err := iw.dc.SavePNG(filepath); err != nil {
			log.WithField("file", filepath).Error("could not save PNG file")
			return 0, err
		}
		*fileNumber++
	}

	return *fileNumber - starting, nil
}

// createTiledFile creates a file with the counter image. Filenumber is the filename, a row and column pointer
// is passed to be able to use the multiplier to create more than one counter in the same sheet
func (iw *imageWriter) createTiledFile(counter *counters.Counter, columns, row *int, canvas *gg.Context) {
	for i := 0; i < counter.Multiplier; i++ {
		canvas.DrawImage(iw.dc.Image(), *columns*counter.Width, *row*counter.Height)
		*columns++
		if *columns == iw.template.Columns {
			*columns = 0
			*row++
		}
	}
	*columns--
}

// GetCanvasForCounter creates a canvas for a counter definition.
func GetCanvasForCounter(cc *counters.Counter, t *counters.CounterTemplate) *gg.Context {
	dc := gg.NewContext(cc.Width, cc.Height)
	if err := dc.LoadFontFace(cc.FontPath, cc.FontHeight); err != nil {
		log.WithField("font", cc.FontPath).Fatal(err)
	}

	dc.Push()
	dc.SetColor(cc.BgColor)
	dc.DrawRectangle(0, 0, float64(cc.Width), float64(cc.Height))
	dc.Fill()
	dc.Pop()

	if cc.FontColorS != "" {
		counters.GetValidColorForString(cc.FontColorS, t.BgColor)
	}

	return dc
}
