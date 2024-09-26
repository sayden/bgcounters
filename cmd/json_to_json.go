package main

import (
	"github.com/pkg/errors"
	"github.com/qdm12/reprint"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/sayden/counters/transform"
	"github.com/thehivecorporation/log"
)

func jsonCountersToJsonFowCounters(counterTemplate *counters.CounterTemplate) (err error) {
	log.Info("creating fow counters")
	t, err := transform.CountersToCounters(
		&transform.CountersToCountersConfig{
			OriginalCounterTemplate: counterTemplate,
			OutputPathInTemplate:    Cli.Json.OutputDestination,
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
		log.WithField("file", Cli.Json.CardTemplateFilepath).WithError(err).Fatal("could not read input file")
	}

	cards, err := transform.CountersToCards(
		&transform.CountersToCardsConfig{
			CountersTemplate: counterTemplate,
			CardTemplate:     cardsTemplate,
			CardCreator: &transform.QuotesToCardBuilder{
				Quotes:         qs,
				IndexForTitles: counterTemplate.IndexNumberForFilename,
			},
		},
	)

	if err != nil {
		return err
	}

	return output.ToJSONFile(cards, Cli.Json.OutputPath)
}

func jsonToBackCounters(counterTemplate *counters.CounterTemplate) (err error) {
	// JSON counters to Back Counters
	finalCounters, err := transform.CountersToCounters(
		&transform.CountersToCountersConfig{
			OriginalCounterTemplate: counterTemplate,
			OutputPathInTemplate:    Cli.Json.OutputDestination,
			CounterBuilder:          &transform.StepLossBackCounterBuilder{},
		},
	)
	if err != nil {
		return errors.Wrap(err, "error trying to convert a counter template into another counter template")
	}

	return output.ToJSONFile(finalCounters, Cli.Json.OutputPath)
}

func jsonPrototypeToJson(counterTemplate *counters.CounterTemplate) (t *counters.CounterTemplate, err error) {
	// JSON counters to Counters, check Prototype in CounterTemplate
	if counterTemplate.Prototypes != nil {
		// ignore counters if prototypes are present
		counterTemplate.Counters = make([]counters.Counter, 0)

		for prototypeName, counter := range counterTemplate.Prototypes {
			log.WithField("prototype", prototypeName).Info("creating counters from prototype")
			// You can prototype texts and images, so one of the two must be present, get their length
			length := 0
			if len(counter.TextsPrototypes) > 0 && len(counter.TextsPrototypes[0].StringList) > 0 {
				length = len(counter.TextsPrototypes[0].StringList)
				if len(counter.ImagesPrototypes) > 0 && len(counter.ImagesPrototypes[0].PathList) > 0 && len(counter.ImagesPrototypes[0].PathList) != length {
					return nil, errors.New("the number of images and texts prototypes must be the same")
				}
			} else if len(counter.ImagesPrototypes) > 0 && len(counter.ImagesPrototypes[0].PathList) > 0 {
				length = len(counter.ImagesPrototypes)
				if len(counter.TextsPrototypes) > 0 && len(counter.TextsPrototypes) != length {
					return nil, errors.New("the number of images and texts prototypes must be the same")
				}
			} else {
				return nil, errors.New("no prototypes found in the counter template")
			}

			for i := 0; i < length; i++ {
				var newCounter counters.Counter
				if err = reprint.FromTo(counter.Counter, &newCounter); err != nil {
					return nil, err
				}

				if counter.TextsPrototypes != nil {
					for _, textPrototype := range counter.TextsPrototypes {
						originalText := counters.Text{}
						if err = reprint.FromTo(textPrototype.Text, &originalText); err != nil {
							return nil, err
						}
						originalText.String = textPrototype.StringList[i]
						newCounter.Texts = append(newCounter.Texts, originalText)
					}
				}

				if counter.ImagesPrototypes != nil {
					for _, imagePrototype := range counter.ImagesPrototypes {
						originalImage := counters.Image{}
						if err = reprint.FromTo(imagePrototype.Image, &originalImage); err != nil {
							return nil, err
						}
						originalImage.Path = imagePrototype.PathList[i]
						newCounter.Images = append(newCounter.Images, originalImage)
					}
				}

				counterTemplate.Counters = append(counterTemplate.Counters, newCounter)
			}
		}

		counterTemplate.Prototypes = nil
		return counterTemplate, nil
	}

	return nil, errors.New("no prototypes found in the counter template")
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
	return output.ToJSONFile(cardTemplate, Cli.Json.OutputDestination)
}