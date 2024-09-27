package counters

import (
	"encoding/json"
	"os"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

type CounterTemplate struct {
	Settings
	Rows    int `json:"rows,omitempty" default:"2"`
	Columns int `json:"columns,omitempty" default:"2"`

	Mode         string  `json:"mode"`
	OutputFolder string  `json:"output_folder" default:"output"`
	DrawGuides   bool    `json:"draw_guides"`
	Scaling      float64 `json:"scaling" default:"1.0"`

	// 0-14 Specify an index number to use when writing a different file for each counter. Multiplier is ignored then
	IndexNumberForFilename int `json:"index_number_for_filename,omitempty" default:"-1"`

	Counters []Counter `json:"counters"`

	Prototypes map[string]CounterPrototype `json:"prototypes,omitempty"`
}

// ParseCountersJsonFile reads a JSON files and parses it into a CounterTemplate after applying it some default
// settings (if not present in the file)
func ParseCountersJsonFile(filepath string, outputFolder ...string) (*CounterTemplate, error) {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read JSON file")
	}

	return ParseTemplate(byt)
}

// ParseTemplate reads a JSON file and parses it into a CounterTemplate after applying it some default settings (if not
// present in the file)
func ParseTemplate(byt []byte) (t *CounterTemplate, err error) {
	t = &CounterTemplate{}
	if err = defaults.Set(t); err != nil {
		return nil, errors.Wrap(err, "could not apply defaults to counter template")
	}

	if err = json.Unmarshal(byt, &t); err != nil {
		log.WithField("incoming_data", string(byt)).Error("could not parse JSON into a counter template")
		return nil, err
	}

	applySettingsScaling(&t.Settings, t.Scaling)

	ApplyCounterWaterfallSettings(t)

	if t.Scaling != 1.0 {
		for i := range t.Counters {
			c := t.Counters[i]
			applyCounterScaling(&c, t.Scaling)
		}
	}

	return
}

func EnrichTemplate(t *CounterTemplate) error {
	if err := defaults.Set(t); err != nil {
		return errors.Wrap(err, "could not read JSON file")
	}

	ApplyCounterWaterfallSettings(t)

	return nil
}

// ApplyCounterWaterfallSettings traverses the counters in the template applying the default settings to value that are
// zero-valued
func ApplyCounterWaterfallSettings(t *CounterTemplate) {
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
	}
}
