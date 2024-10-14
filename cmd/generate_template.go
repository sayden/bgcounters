package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/xeipuuv/gojsonschema"
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
			FontHeight:      10,
			Width:           100,
			Height:          100,
			Margins:         3,
			FontHeight:      10,
			FontColorS:      "black",
			BackgroundColor: "white",
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

func validateSchema(inputPath string) error {
	logger.Info("Validating JSON file")

	r := new(jsonschema.Reflector)
	counterTemplateSchemaMarshaller := r.Reflect(&counters.CounterTemplate{})
	byt, err := counterTemplateSchemaMarshaller.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "could not marshal counter template schema")
	}

	schema := gojsonschema.NewBytesLoader(byt)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + inputPath)
	result, err := gojsonschema.Validate(schema, documentLoader)
	if err != nil {
		return errors.Wrap(err, "could not validate JSON file")
	}

	if !result.Valid() {
		logger.Error("The document is not valid. see errors: ")
		for _, desc := range result.Errors() {
			logger.Errorf("- %s", desc)
		}
		return errors.Wrap(err, "JSON file is not valid")
	} else {
		logger.Debug("JSON file is valid")
	}

	return nil
}
