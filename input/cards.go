package input

import (
	"encoding/json"
	"io"
	"os"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

func ReadJSONCardsFile(cardsFilepath string) (*counters.CardsTemplate, error) {
	f, err := os.Open(cardsFilepath)
	if err != nil {
		log.WithField("file", cardsFilepath).Error("could not open cards file")
		return nil, err
	}
	defer f.Close()

	return parseCardsTemplateJSON(f)
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
