package main

import (
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/sayden/counters/transform"
)

func jsonCountersToJsonFowCounters(counterTemplate *counters.CounterTemplate) (err error) {
	logger.Info("creating fow counters")

	t, err := transform.CountersToCounters(
		&transform.CountersToCountersConfig{
			OriginalCounterTemplate: counterTemplate,
			OutputPathInTemplate:    Cli.Json.Destination,
			CounterBuilder:          &transform.SimpleFowCounterBuilder{},
		},
	)
	if err != nil {
		return errors.Wrap(err, "error trying to convert a counter template into another counter template")
	}

	return output.ToJSONFile(t, Cli.Json.OutputPath)
}

func jsonCountersToJsonCards(counterTemplate *counters.CounterTemplate) (err error) {
	qs, err := input.ReadQuotesFromFile(Cli.Json.QuotesFile)
	if err != nil {
		return errors.Wrap(err, "could not read quotes file")
	}

	if Cli.Json.CardTemplateFilepath == "" {
		return errors.New("A card template must be provided using 'card-template-filepath' when writing a card output")
	}
	cardsTemplate, err := input.ReadJSONCardsFile(Cli.Json.CardTemplateFilepath)
	if err != nil {
		logger.Fatal("could not read input file", "file", Cli.Json.CardTemplateFilepath, "error", err)
	}

	cards, err := transform.CountersToCards(
		&transform.CountersToCardsConfig{
			CountersTemplate: counterTemplate,
			CardTemplate:     cardsTemplate,
			CardCreator: &transform.QuotesToCardBuilder{
				Quotes:         qs,
				IndexForTitles: counterTemplate.PositionNumberForFilename,
			},
		},
	)

	if err != nil {
		return err
	}

	return output.ToJSONFile(cards, Cli.Json.OutputPath)
}

func jsonToJsonCardEvents(events []counters.Event) (err error) {
	images, err := input.GetFilenamesForPath(Cli.Json.BackgroundImages)
	if err != nil {
		return errors.Wrap(err, "error trying to load bg images")
	}

	cardTemplate := transform.EventsToCards(
		&transform.EventsToCardsConfig{
			Events:             events,
			Images:             images,
			BackImageFile:      Cli.Json.BackImage,
			GeneratedImageName: Cli.Json.OutputPath,
		},
	)
	return output.ToJSONFile(cardTemplate, Cli.Json.Destination)
}
