package transform

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type CountersToCardsConfig struct {
	CountersTemplate *counters.CounterTemplate
	CardTemplate     *counters.CardsTemplate
	CardCreator      counters.CardBuilder
}

// CountersToCards using the provided input writes png images into files in batches of 70 (70 cards per file).
func CountersToCards(cfg *CountersToCardsConfig) (*counters.CardsTemplate, error) {
	var outputCardTemplate counters.CardsTemplate
	mergo.Merge(&outputCardTemplate, &cfg.CardTemplate)

	// INPUT
	for _, counter := range cfg.CountersTemplate.Counters {
		if !counter.Extra.SkipCardGeneration {
			// convert counter to fow card
			card, err := cfg.CardCreator.ToCard(counter, cfg.CardTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "error trying to create card")
			}

			if card == nil {
				return nil, errors.New("error creating card")
			}

			for i := 0; i < counter.Multiplier; i++ {
				outputCardTemplate.Cards = append(outputCardTemplate.Cards, *card)
			}
		}
	}

	return &outputCardTemplate, nil
}
