// Description: This file contains the logic to generate images from the JSON files.
// It uses the counters package to generate the images. The AssetsOutput struct is used to define
// the CLI options for the assets command.
package main

import (
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/thehivecorporation/log"
)

const (
	OutputContentCounters = iota
	OutputContentBackCounters
	OutputContentCards
	OutputContentFowCounters
	OutputContentEvents
	OutputContentBlock
)

const (
	AssetFormatTile = iota
	AssetFormatIndividual
)

type AssetsOutput struct {
	OutputContent string `help:"InputContent to produce: counters, blocks or cards" required:"true"`
	InputPath     string `help:"Input path of the file to read. Be aware that some outputs requires specific inputs." short:"i" required:"true"`
	OutputPath    string `help:"Path to the folder to write the image(s)" short:"o"`
	Tiled         bool   `help:"Write a sheet of 7x10 items per parge" default:"false"`
	Individual    bool   `help:"Write a file for each counter/card" default:"true"`
	BlockBack     string `help:"If using --output-content blocks, set the input path of the JSON to place in the back of the counter if it applies"`
}

func (i *AssetsOutput) Run(ctx *kong.Context) error {
	switch Cli.Assets.OutputContent {
	// Outputs blocks images
	case "blocks":
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
		if i.BlockBack != "" {
			if backCounterTemplate, err = input.ReadCounterTemplate(Cli.Assets.BlockBack); err != nil {
				return err
			}
		}

		return output.CountersToBlocks(counterTemplate, backCounterTemplate)
	case "counters":
		counterTemplate, err := counters.ParseCountersJsonFile(Cli.Assets.InputPath, Cli.Assets.OutputPath)
		if err != nil {
			return err
		}

		// Override output path with the one provided in the CLI
		if Cli.Assets.OutputPath != "" {
			log.WithField("output_path", Cli.Assets.OutputPath).Info("Overriding output path")
			counterTemplate.OutputFolder = Cli.Assets.OutputPath
		}

		return output.CountersToPNG(counterTemplate)
	case "cards":
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

	return errors.New("'--output-content' not recognized for Image output")
}
