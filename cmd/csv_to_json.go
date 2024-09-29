package main

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/thehivecorporation/log"
)

func csvToCounters(byt []byte) (err error) {
	if Cli.Json.CounterTemplateFilepath == "" {
		return errors.New("a counter template must be provided using 'card-template-filepath'")
	}

	if Cli.Json.OutputPath == "" {
		return errors.New("an output path for the output file is required")
	}

	countersTemplate, err := input.ReadCounterTemplate(Cli.Json.CounterTemplateFilepath)
	if err != nil {
		return err
	}

	csvCounters, err := input.ReadCSVCounters(Cli.Json.InputPath)
	if err != nil {
		return err
	}

	if err = mergo.Merge(&csvCounters, countersTemplate); err != nil {
		return errors.Wrap(err, "could not merge CSV counters with provided template")
	}

	return output.ToJSONFile(csvCounters, Cli.Json.OutputPath)
}

func csvToCards(byt []byte) (err error) {
	if Cli.Json.CardTemplateFilepath == "" {
		return errors.New("A card template must be provided using 'card-template-filepath'")
	}

	cardsTemplate, err := input.ReadJSONCardsFile(Cli.Json.CardTemplateFilepath)
	if err != nil {
		log.WithField("file", Cli.Json.CardTemplateFilepath).WithError(err).Fatal("could not read card-template-filepath")
	}

	// Generate a JSON version of the incoming CSV
	content, err := input.ReadCSVCardsFromBytes(byt, cardsTemplate)
	if err != nil {
		return err
	}

	if Cli.Json.Destination != "" {
		content.OutputPath = Cli.Json.Destination
	}
	content.OutputPath = Cli.Json.Destination

	if Cli.Json.OutputPath == "" {
		return errors.New("an output path for the output file is required")
	}

	return output.ToJSONFile(content, Cli.Json.OutputPath)

}
