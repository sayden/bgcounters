package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/sayden/counters/transform"
	"github.com/thehivecorporation/log"
)

type JsonOutput struct {
	InputPath  string `help:"Input path of the file to read" short:"i" required:"true"`
	OutputPath string `help:"Path to the folder to write the JSON" short:"o"`
	OutputType string `help:"Type of content to produce: back-counters, cards, fow-counters, counters or events" short:"c"`

	EventsPoolFile          string `help:"A file to take 'events' from"`
	BackImage               string `help:"The image for the back of the cards"`
	OutputDestination       string `help:"When generating a JSON Template, this contains the destination folder for images inside the template"`
	CardTemplateFilepath    string `help:"When writing cards from a CSV, a template for those cards must be provided"`
	CounterTemplateFilepath string `help:"When writing counters from a CSV, a template for those counters must be provided"`
	BackgroundImages        string `help:"Path to a folder containing background images for the cards"`
	QuotesFile              string `help:"Path to a JSON file containing quotes for the cards"`
}

func (i *JsonOutput) Run(ctx *kong.Context) error {
	var counterTemplate *counters.CounterTemplate
	var events []counters.Event
	var err error

	// content could be JSON or CSV
	inputContent, err := fsops.GetExtension(Cli.Json.InputPath)
	if err != nil {
		return err
	}

	byt, err := os.ReadFile(Cli.Json.InputPath)
	if err != nil {
		return errors.Wrap(err, "could not read input file")
	}

	switch inputContent {
	case counters.FileContent_CSV:
		switch Cli.Json.OutputType {
		case "cards":
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

			if Cli.Json.OutputDestination != "" {
				content.OutputPath = Cli.Json.OutputDestination
			}
			content.OutputPath = Cli.Json.OutputDestination

			if Cli.Json.OutputPath == "" {
				return errors.New("an output path for the output file is required")
			}

			return output.ToJSONFile(content, Cli.Json.OutputPath)
		case "counters":
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
	case counters.FileContent_JSON:
		inputContent, err := fsops.IdentifyJSONFileContent(byt)
		if err != nil {
			return errors.Wrap(err, "could not identify file content")
		}

		switch inputContent {
		case counters.FileContent_CounterTemplate:
			counterTemplate, err = input.ReadCounterTemplate(Cli.Json.InputPath, Cli.Json.OutputPath)
			if err != nil {
				return errors.Wrap(err, "error reading counter template")
			}
		case counters.FileContent_Events:
			events, err = input.JsonEventsToEvents(Cli.Json.InputPath)
			if err != nil {
				return errors.Wrap(err, "could not read events file")
			}
		default:
			return errors.New("combination of options not recognized")
		}

		switch Cli.Json.OutputType {
		case "counters":
			// JSON counters to Counters, check Prototype in CounterTemplate
			if counterTemplate.Prototypes != nil {
				// ignore counters if prototypes are present
				counterTemplate.Counters = make([]counters.Counter, 0)

				for _, counter := range counterTemplate.Prototypes {
					multiplier := counter.Multiplier
					if multiplier == 0 {
						multiplier = 1
					}
					counter.Multiplier = 0

					for i := 0; i < multiplier; i++ {
						counterTemplate.Counters = append(counterTemplate.Counters, counter)
					}
				}

				templ, err := counters.ParseTemplate(byt)
				if err != nil {
					return errors.Wrap(err, "could not parse counter template")
				}

				err = output.ToJSONFile(templ, Cli.Json.OutputPath)

				return err
			}

			return errors.New("no prototypes found in the counter template")
		case "back-counters":
			// JSON counters to Back Counters
			finalCounters, err := transform.CountersToCounters(
				&transform.CountersToCountersConfig{
					OriginalCounterTemplate: counterTemplate,
					CounterBuilder:          &transform.SimpleFowCounterBuilder{},
				},
			)
			if err != nil {
				return errors.Wrap(err, "error trying to convert a counter template into another counter template")
			}

			return output.ToJSONFile(finalCounters, Cli.Json.OutputPath)
		case "cards":
			// JSON counters to Cards
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
		case "fow-counters":
			// JSON counters to Fow Counters
			log.Info("creating fow counters")
			t, err := transform.CountersToCounters(
				&transform.CountersToCountersConfig{
					OriginalCounterTemplate: counterTemplate,
					OutputPathInTemplate:    Cli.Json.OutputDestination,
					CounterBuilder:          &transform.SimpleFowCounterBuilder{},
				},
			)
			ctx.FatalIfErrorf(err)
			return output.ToJSONFile(t, Cli.Json.OutputPath)
		case "events":
			// FIXME JSON to Events
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
	}

	return errors.New("'--output-content' not recognized for JSON output")
}
