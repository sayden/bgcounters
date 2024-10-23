package input

import (
	"encoding/json"
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

	var t counters.CardsTemplate
	if err := defaults.Set(&t.Settings); err != nil {
		return nil, errors.Wrap(err, "could not set defaults into card template")
	}

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&t); err != nil {
		return nil, errors.Wrap(err, "could not read JSON card data")
	}

	err = t.ApplyCardWaterfallSettings()
	if err != nil {
		return nil, errors.Wrap(err, "could not apply card waterfall settings")
	}

	return &t, nil
}
