package output

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

func CardsToPNG(template *counters.CardsTemplate) error {
	sheetFileNumber := 1

	// count the total amounts to card that will be processed
	totalCards := 0
	for _, card := range template.Cards {
		totalCards += card.Multiplier
	}

	var (
		rows, columns int
		cards         []*counters.Card
		currentIndex  int
		n             int
	)

	log.Infof("Generating a total of %d cards in sheets of %d cards", totalCards, template.Rows*template.Columns)

	for {
		sheet, err := createOutputCanvas(template)
		if err != nil {
			return errors.Wrap(err, "could not create canvas to write cards")
		}

		cardProc := newCardProcessor(&cardProcessorConfig{
			template:   template,
			cardCanvas: sheet,
		})

	processing:
		for row := 0; row < rows; row++ {
			for col := 0; col < columns; col++ {
				if currentIndex >= len(cards) {
					break processing
				}

				card := cards[currentIndex]
				if err := cardProc.processCard(sheet, card, col, row); err != nil {
					return errors.Wrap(err, "error trying to process card")
				}
				currentIndex++
			}
		}

		outputPath := fmt.Sprintf(template.OutputPath, sheetFileNumber)

		if n >= template.Rows*template.Columns {
			// More than one file required to write cards. Expect a Go's format string in the output path
			sheetFileNumber++
		}

		// Save result on a file
		if err = sheet.SavePNG(outputPath); err != nil {
			log.WithField("output_path", outputPath).Error("error trying to save final cards image file")
			return err
		} else {
			log.WithFields(log.Fields{"output_file_path": outputPath, "generated_cards": n}).Info("Written card file")
		}

		if n >= totalCards {
			break
		}
	}

	return nil
}

func createOutputCanvas(template *counters.CardsTemplate) (*gg.Context, error) {
	width := template.Columns * template.Width
	height := template.Rows * template.Height
	sheet := gg.NewContext(width, height)
	if err := sheet.LoadFontFace(counters.DEFAULT_FONT_PATH, template.FontHeight); err != nil {
		log.WithField("font_path", counters.DEFAULT_FONT_PATH).Error("could not load font face")
		return nil, err
	}

	sheet.SetColor(color.White)
	sheet.Fill()

	return sheet, nil
}
