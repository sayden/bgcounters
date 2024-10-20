package counters

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
)

type CounterTemplate struct {
	Settings

	WorkingDirectory string `json:"working_directory,omitempty"`

	Rows    int `json:"rows,omitempty" default:"2" jsonschema_description:"Number of rows, required when creating tiled based sheets for printing or TTS"`
	Columns int `json:"columns,omitempty" default:"2" jsonschema_description:"Number of columns, required when creating tiled based sheets for printing or TTS"`

	Mode         string  `json:"mode"`
	OutputFolder string  `json:"output_folder" default:"output"`
	DrawGuides   bool    `json:"draw_guides,omitempty"`
	Scaling      float64 `json:"scaling,omitempty" default:"1.0"`

	// 0-16 Specify an position in the counter to use when writing a different file
	PositionNumberForFilename int `json:"position_number_for_filename,omitempty"`

	Counters []Counter `json:"counters,omitempty"`

	Prototypes map[string]CounterPrototype `json:"prototypes,omitempty"`
}

// ParseCounterTemplate reads a JSON file and parses it into a CounterTemplate after applying it some default settings (if not
// present in the file)
func ParseCounterTemplate(byt []byte) (t *CounterTemplate, err error) {
	if err = ValidateSchemaBytes[CounterTemplate](byt); err != nil {
		return nil, errors.Wrap(err, "JSON file is not valid")
	}

	t = &CounterTemplate{}
	if err = defaults.Set(t); err != nil {
		return nil, errors.Wrap(err, "could not apply defaults to counter template")
	}

	if err = json.Unmarshal(byt, &t); err != nil {
		return nil, err
	}

	t.Settings.ApplySettingsScaling(t.Scaling)

	t.ApplyCounterWaterfallSettings()

	// Request body contains the current working directory to use
	// This is relevant because we need to use relavite paths
	if t.WorkingDirectory != "" {
		if err = os.Chdir(os.ExpandEnv(t.WorkingDirectory)); err != nil {
			return nil, err
		}
	}

	if t.Scaling != 1.0 {
		for i := range t.Counters {
			c := t.Counters[i]
			ApplyCounterScaling(&c, t.Scaling)
		}
	}

	return
}

func (t *CounterTemplate) EnrichTemplate() error {
	if err := defaults.Set(t); err != nil {
		return errors.Wrap(err, "could not read JSON file")
	}

	t.ApplyCounterWaterfallSettings()

	return nil
}

// ApplyCounterWaterfallSettings traverses the counters in the template applying the default settings to value that are
// zero-valued
func (t *CounterTemplate) ApplyCounterWaterfallSettings() {
	SetColors(&t.Settings)

	for counterIndex, counter := range t.Counters {
		Merge(&t.Counters[counterIndex].Settings, t.Settings)
		if t.Counters[counterIndex].Back != nil {
			Merge(&t.Counters[counterIndex].Back.Settings, t.Settings)
		}

		for imageIndex := range counter.Images {
			Merge(&t.Counters[counterIndex].Images[imageIndex].Settings, counter.Settings)
			if t.Counters[counterIndex].Back != nil {
				Merge(&t.Counters[counterIndex].Back.Images[imageIndex].Settings, t.Settings)
			}
		}

		for imageIndex := range counter.Texts {
			Merge(&t.Counters[counterIndex].Texts[imageIndex].Settings, counter.Settings)
			if t.Counters[counterIndex].Back != nil {
				Merge(&t.Counters[counterIndex].Back.Texts[imageIndex].Settings, t.Settings)
			}
		}

		if counter.Multiplier == 0 {
			counter.Multiplier = 1
		}
	}
}

func (ct *CounterTemplate) ParsePrototype() (*CounterTemplate, error) {
	// JSON counters to Counters
	newTemplate, err := ct.ExpandPrototypeCounterTemplate()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to convert a counter template into another counter template")
	}
	newTemplate.Scaling = 1

	byt, err := json.Marshal(newTemplate)
	if err != nil {
		return nil, err
	}

	newTemplate, err = ParseCounterTemplate(byt)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse JSON file")
	}

	return newTemplate, nil
}

func (ct *CounterTemplate) ExpandPrototypeCounterTemplate() (t *CounterTemplate, err error) {
	// JSON counters to Counters, check Prototype in CounterTemplate
	if ct.Prototypes != nil {
		if ct.Counters == nil {
			ct.Counters = make([]Counter, 0)
		}

		// sort prototypes by name, to ensure consistent output filenames this is a small
		// inconvenience, because iterating over maps in Go returns keys in random order
		names := make([]string, 0, len(ct.Prototypes))
		for name := range ct.Prototypes {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, prototypeName := range names {
			prototype := ct.Prototypes[prototypeName]

			cts, err := prototype.ToCounters()
			if err != nil {
				return nil, err
			}

			ct.Counters = append(ct.Counters, cts...)
		}

		ct.Prototypes = nil
		return ct, nil
	}

	return ct, nil
}
