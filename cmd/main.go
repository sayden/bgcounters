package main

import (
	"github.com/alecthomas/kong"
	"github.com/thehivecorporation/log"
)

var Cli struct {
	// Asset option is used to generate assets directly. This is usually counters images or Card sheets.
	Assets AssetsOutput `cmd:"" help:"Generate images of some short, using either counters or cards, from a JSON file"`

	// JSON uses a JSON input to generate another JSON output.
	Json JsonOutput `cmd:"" help:"Generate a JSON of some short, by transforming another JSON as input"`

	// Vassal is used to generate a Vassal module for testing purposes.
	Vassal vassal `cmd:"" help:"Create a vassal module for testing. It searches for the 'template.xml' in the same folder"` //FIXME
}

func main() {
	ctx := kong.Parse(&Cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)

	log.Info("Done!")
}
