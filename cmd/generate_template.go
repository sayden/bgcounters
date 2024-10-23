package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type GenerateTemplate struct {
	OutputPath string `help:"Path to the JSON output file" short:"o"`
}

func (i *GenerateTemplate) Run(ctx *kong.Context) error {
	if i.OutputPath == "" {
		return errors.New("output path is required")
	}

	err := generateNewCounterTemplate(i.OutputPath)
	if err != nil {
		return errors.Wrap(err, "could not generate new counter template")
	}

	return nil
}

func generateNewCounterTemplate(outputPath string) error {
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
				TextPrototypes: []counters.TextPrototype{
					{
						StringList: []string{"String1", "String2"},
					},
				},
				ImagePrototypes: []counters.ImagePrototype{
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
			FontHeight:      10,
			Width:           100,
			Height:          100,
			Margins:         floatP(3),
			FontColorS:      "black",
			BackgroundColor: stringP("white"),
			ImageScaling:    "fitWidth",
		},
	}

	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not create output directories")
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "could not create output file")
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(counterTemplate)
}
