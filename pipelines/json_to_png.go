package pipelines

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/thehivecorporation/log"
)

// JSONCardsToPNG output is always to write a single PNG file with the cards. It requires a valid path to a JSON in
// an expected format.
func JSONCardsToPNG(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	template, err := input.ParseCardsFileTemplate(f)
	if err != nil {
		return errors.Wrap(err, "error trying to parse cards file")
	}

	log.WithFields(log.Fields{"source_filepath": filepath, "cards": len(template.Cards)}).Info("Generating cards")
	return output.CardsToPNG(template)
}

// JSONCountersToPNG generate one single cardCanvas PNG image of counters or many PNG images with counters.
// it requires a valid `filepath` to a JSON in an expected format.
func JSONCountersToPNG(filepath string) error {
	log.WithField("filepath", filepath).Info("Generating counters")

	byt, err := os.ReadFile(filepath)
	if err != nil {
		log.WithField("filepath", filepath).Error(err)
		return err
	}

	template, err := counters.ParseTemplate(byt)
	if err != nil {
		return errors.Wrap(err, "could not parse template")
	}

	if err = output.CountersToPNG(template); err != nil {
		return errors.Wrap(err, "could not write PNG for counters")
	}

	return nil
}
