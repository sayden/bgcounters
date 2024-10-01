// Description: This file contains the logic to generate images from the JSON files.
// It uses the counters package to generate the images. The AssetsOutput struct is used to define
// the CLI options for the assets command.
package main

import (
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
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
	OutputType string `help:"InputContent to produce: counters, blocks or cards" required:"true" short:"t" default:"counters"`
	InputPath  string `help:"Input path of the file to read. Be aware that some outputs requires specific inputs." short:"i" required:"true"`
	OutputPath string `help:"Path to the folder to write the image(s)" short:"o" default:"./generated"`

	Tiled      bool   `help:"Write a sheet of 7x10 items per parge" default:"false"`
	Individual bool   `help:"Write a file for each counter/card" default:"true"`
	BlockBack  string `help:"If using --output-content blocks, set the input path of the JSON to place in the back of the counter if it applies"`
}

func (i *AssetsOutput) Run(ctx *kong.Context) error {
	switch Cli.Assets.OutputType {
	// Outputs blocks images
	case "blocks":
		return jsonToBlock(i.BlockBack)
	case "counters":
		return jsonToAsset(Cli.Assets.InputPath, Cli.Assets.OutputPath)
	case "cards":
		return jsonToCards()
	}

	return errors.New("'--output-content' not recognized for Image output")
}
