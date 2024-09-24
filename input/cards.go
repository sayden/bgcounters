package input

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

// ParseCardsFileTemplate reads 'filepath' and returns a parsed CardsTemplate from it
func ParseCardsFileTemplate(f io.Reader) (*counters.CardsTemplate, error) {
	return parseCardsTemplateJSON(f)
}

func ReadJSONCardsFile(cardsFilepath string) (*counters.CardsTemplate, error) {
	f, err := os.Open(cardsFilepath)
	if err != nil {
		log.WithField("file", cardsFilepath).Error("could not open cards file")
		return nil, err
	}
	defer f.Close()

	return ParseCardsFileTemplate(f)
}

func ReadCSVCardsFromBytes(byt []byte, template *counters.CardsTemplate) (*counters.CardsTemplate, error) {
	return ReadCSVCards(bytes.NewReader(byt), template)
}

func ReadCSVCards(f io.Reader, template *counters.CardsTemplate) (*counters.CardsTemplate, error) {
	csvReader := csv.NewReader(f)

	headTitles, err := csvReader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to read the first row on the csv file")
	}

	cards := make([]counters.Card, 0)

	for {
		cols, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		card := counters.Card{
			Areas: make([]counters.Counter, 0),
		}

		// skip leftmost column (multiplier) by now
		for i, col := range cols[1:] {
			if i == 0 && headTitles[1] == "bg_color" {
				card.Settings.BackgroundColor = col
				continue
			}

			cell := strings.TrimSpace(col)

			card.Areas = append(card.Areas, counters.Counter{
				Texts: []counters.Text{{
					Settings: counters.Settings{
						Position:      0,
						AvoidClipping: true,
						BgColor:       color.White,
					},
					String: cell,
				}},
			})
		}

		// read the leftmost column now
		multiplier, err := strconv.Atoi(cols[0])
		if err != nil {
			return nil, err
		}

		for i := 0; i < multiplier; i++ {
			cards = append(cards, card)
		}
	}

	template.Cards = cards

	if err := defaults.Set(&template.Settings); err != nil {
		return nil, errors.Wrap(err, "could not set defaults into card template")
	}

	// counters.ApplyCardWaterfallSettings(template)

	// if err = counters.EnrichTemplate(template); err != nil {
	// 	return nil, errors.Wrap(err, "could not enrich template")
	// }

	return template, nil
}

// parseCardsTemplateJSON reads the input file, parses the JSON, sets defaults for settings and apply the known settings
// in waterfall effect to every card in the list.
// This is useful when generating a JSON CardTemplate file from counters, because when Marshalling the content of the
// contents to JSON, all the settings on the counters are written even if they are empty.
func parseCardsTemplateJSON(f io.Reader) (*counters.CardsTemplate, error) {
	var t counters.CardsTemplate
	if err := defaults.Set(&t.Settings); err != nil {
		return nil, errors.Wrap(err, "could not set defaults into card template")
	}

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&t); err != nil {
		return nil, errors.Wrap(err, "could not read JSON card data")
	}

	counters.ApplyCardWaterfallSettings(&t)

	return &t, nil
}
