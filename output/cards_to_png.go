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
	proc := newSheetProcessor(template)

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

		n, err := proc.processAll(sheet, cardProc)
		if err != nil {
			return err
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

type sheetProcessor struct {
	rows, columns int
	cards         []*counters.Card
	currentIndex  int
}

func newSheetProcessor(template *counters.CardsTemplate) *sheetProcessor {
	cards := make([]*counters.Card, 0, len(template.Cards))
	for _, card := range template.Cards {
		for i := 0; i < card.Multiplier; i++ {
			c := card
			cards = append(cards, &c)
		}
	}

	return &sheetProcessor{
		rows:    template.Rows,
		columns: template.Columns,
		cards:   cards,
	}
}

func (s *sheetProcessor) processAll(sheet *gg.Context, processor *cardProcessor) (int, error) {
	for row := 0; row < s.rows; row++ {
		for col := 0; col < s.columns; col++ {
			if s.currentIndex >= len(s.cards) {
				return s.currentIndex, nil
			}

			card := s.cards[s.currentIndex]
			if err := processor.processCard(sheet, card, col, row); err != nil {
				return 0, errors.Wrap(err, "error trying to process card")
			}
			s.currentIndex++
		}
	}

	return s.currentIndex, nil
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

func drawFrame(c *gg.Context, w float64, col color.Color) {
	c.Push()
	c.SetColor(col)
	c.SetLineWidth(w)
	frameX := float64(c.Width())
	frameY := float64(c.Height())
	c.DrawRectangle(0, 0, frameX, frameY)
	c.Stroke()
	c.Pop()
}
