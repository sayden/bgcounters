package main

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

var logger = log.New(os.Stderr)

var Cli struct {
	// Asset option is used to generate assets directly. This is usually counters images or Card sheets.
	Assets AssetsOutput `cmd:"" help:"Generate images of some short, using either counters or cards, from a JSON file"`

	// JSON uses a JSON input to generate another JSON output.
	Json JsonOutput `cmd:"" help:"Generate a JSON of some short, by transforming another JSON as input"`

	// Vassal is used to generate a Vassal module for testing purposes.
	Vassal vassal `cmd:"" help:"Create a vassal module for testing. It searches for the 'template.xml' in the same folder"` //FIXME

	New NewDefaultTemplates `cmd:"" help:"Generates a new counter template file with default values"`
}

func main() {
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)

	ctx := kong.Parse(&Cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)

	log.Info("Done")
}

type NewDefaultTemplates struct {
	OutputPath string `help:"Path to the folder to write the JSON" short:"o"`
}

func (i *NewDefaultTemplates) Run(ctx *kong.Context) error {
	if i.OutputPath == "" {
		return errors.New("output path is required")
	}

	err := generateNewCounterTemplate(i.OutputPath)
	if err != nil {
		return errors.Wrap(err, "could not generate new counter template")
	}

	return nil
}

func generateNewCounterTemplate(outputhPath string) error {
	counterTemplate := counters.CounterTemplate{
		Counters: []counters.Counter{
			{
				Texts: []counters.Text{
					{
						Settings: counters.Settings{
							Position: 3,
						},
						String: "String",
					},
				},
				Images: []counters.Image{
					{
						Settings: counters.Settings{
							Position: 0,
						},
						Path: "path/to/image",
					},
				},
			},
		},
		Rows:                      7,
		Columns:                   7,
		Mode:                      "template",
		OutputFolder:              "generated",
		PositionNumberForFilename: 0,
		Prototypes: map[string]counters.CounterPrototype{
			"prototype1": {
				TextsPrototypes: []counters.TextPrototype{
					{
						StringList: []string{"String1", "String2"},
					},
				},
				ImagesPrototypes: []counters.ImagePrototype{
					{
						Image: counters.Image{
							Settings: counters.Settings{
								Position: 0,
							},
						},
						PathList: []string{"path/to/image1", "path/to/image2"},
					},
				},
				Counter: counters.Counter{
					Texts: []counters.Text{
						{
							Settings: counters.Settings{
								Position: 3,
							},
							String: "String",
						},
					},
					Images: []counters.Image{
						{
							Settings: counters.Settings{
								Position: 0,
							},
							Path: "path/to/image",
						},
					},
				},
			},
		},
		Settings: counters.Settings{
			FontPath:        "BebasNeue-Regular.ttf",
			Width:           100,
			Height:          100,
			Margins:         3,
			FontHeight:      10,
			FontColorS:      "black",
			BackgroundColor: "white",
		},
	}

	f, err := os.Create(outputhPath)
	if err != nil {
		return errors.Wrap(err, "could not create output file")
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(counterTemplate)
}
