package input

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

// ReadCounterTemplate is used to read files the defines an array of counters. `outputPathForTemplate` is only used
// for Cli stuff
func ReadCounterTemplate(inputPath string, outputPathForTemplate ...string) (*counters.CounterTemplate, error) {
	extension := filepath.Ext(inputPath)

	switch extension {
	case ".csv":
		counterTemplate, err := ReadCSVCounters(inputPath)
		if err != nil {
			return nil, err
		}

		if len(outputPathForTemplate) == 1 {
			counterTemplate.OutputFolder = outputPathForTemplate[0]
		}

		return counterTemplate, nil
	case ".json":
		return counters.ParseCountersJsonFile(inputPath, outputPathForTemplate...)
	}

	return nil, fmt.Errorf("extension '%s' not found", extension)
}

// ReadCSVCounters follows a convention to read CSV files into counter templates:
// A-Q Columns in the CSV corresponds to the 0-16 possible position to write into a counter, being 0 the
// horizontal and vertical center
// Column R (17): Is the Side: red, blue, german, etc.
// Column S (18): Is the background color
func ReadCSVCounters(filepath string) (*counters.CounterTemplate, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to open counter template")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	// Skip headers
	_, err = csvReader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to read the first row on the csv file")
	}

	cnts := make([]counters.Counter, 0)

	for {
		cols, err := csvReader.Read()
		if err != nil {
			break
		}

		cnt := counters.Counter{
			Texts: make([]counters.Text, 0),
		}

		for colIndex, col := range cols {
			if col != "" {
				cell := strings.TrimSpace(col)

				// Side: red, blue, german, russian, etc.
				if colIndex == 17 {
					cnt.Extra.Side = cell
					continue
				} else if colIndex == 18 {
					//background color
					cnt.BackgroundColor = "#" + cell
					break
				}

				cnt.Texts = append(
					cnt.Texts, counters.Text{
						Settings: counters.Settings{
							Position:      colIndex,
							AvoidClipping: true,
							FontColorS:    "white",
							StrokeWidth:   2,
							StrokeColorS:  "black",
						},
						String: cell,
					},
				)
			}
		}

		cnts = append(cnts, cnt)
	}

	//TODO remove hardcoded data
	template := counters.CounterTemplate{
		Rows:                         7,
		Columns:                      7,
		Mode:                         "template",
		OutputFolder:                 "/tmp",
		Counters:                     cnts,
		PositionNumberForFilename: 3,
	}

	template.Settings.Width = 50
	template.Settings.Height = 50
	template.Settings.FontHeight = 16
	template.Settings.Margins = 2
	template.Settings.FontColor = color.White
	template.Settings.ImageScaling = counters.IMAGE_SCALING_FIT_NONE

	if err = counters.EnrichTemplate(&template); err != nil {
		return nil, errors.Wrap(err, "could not enrich template")
	}

	return &template, nil
}
