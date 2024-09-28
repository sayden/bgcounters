package main

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/thehivecorporation/log"
)

func jsonToCards() (err error) {
	cardsTemplate, err := input.ReadJSONCardsFile(Cli.Assets.InputPath)
	if err != nil {
		return errors.Wrap(err, "error trying to read json cards file")
	}

	// Override output path with the one provided in the CLI
	if Cli.Assets.OutputPath != "" {
		log.WithField("output_path", Cli.Assets.OutputPath).Info("Overriding output path")
		cardsTemplate.OutputPath = Cli.Assets.OutputPath
	}

	return output.CardsToPNG(cardsTemplate)
}

func jsonToAsset(inputPath, outputPath string) (err error) {
	counterTemplate, err := input.ReadCounterTemplate(inputPath, outputPath)
	if err != nil {
		return errors.Wrap(err, "error reading counter template")
	}

	// JSON counters to Counters
	newTemplate, err := jsonPrototypeToJson(counterTemplate)
	if err != nil {
		return errors.Wrap(err, "error trying to convert a counter template into another counter template")
	}
	newTemplate.Scaling = 1

	byt, err := json.Marshal(newTemplate)
	if err != nil {
		return err
	}

	newTemplate, err = counters.ParseTemplate(byt)
	if err != nil {
		return errors.Wrap(err, "could not parse JSON file")
	}

	// Override output path with the one provided in the CLI
	if outputPath != "" {
		log.WithField("output_path", outputPath).Info("Overriding output path")
		newTemplate.OutputFolder = outputPath
	}

	return output.CountersToPNG(newTemplate)
}

func jsonToBlock(blockBack string) (err error) {
	counterTemplate, err := input.ReadCounterTemplate(Cli.Assets.InputPath, Cli.Assets.OutputPath)
	if err != nil {
		return err
	}

	// Override output path with the one provided in the CLI
	if Cli.Assets.OutputPath != "" {
		log.WithField("output_path", Cli.Assets.OutputPath).Info("Overriding output path")
		counterTemplate.OutputFolder = Cli.Assets.OutputPath
	}

	var backCounterTemplate *counters.CounterTemplate
	if blockBack != "" {
		if backCounterTemplate, err = input.ReadCounterTemplate(Cli.Assets.BlockBack); err != nil {
			return err
		}
	}

	return output.CountersToBlocks(counterTemplate, backCounterTemplate)

}
