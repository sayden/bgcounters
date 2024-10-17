package input

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

// ReadCounterTemplate is used to read files the defines an array of counters. `outputPathForTemplate` is only used
// for Cli stuff
func ReadCounterTemplate(inputPath string, outputPathForTemplate ...string) (*counters.CounterTemplate, error) {
	extension := filepath.Ext(inputPath)

	switch extension {
	case ".csv":
		counterTemplate, err := ReadCSVCounters(inputPath)
		if err != nil {
			return nil, err
		}

		if len(outputPathForTemplate) == 1 {
			counterTemplate.OutputFolder = outputPathForTemplate[0]
		}

		return counterTemplate, nil
	case ".json":
		byt, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not read JSON file")
		}

		return counters.ParseCounterTemplate(byt)
	}

	return nil, fmt.Errorf("extension '%s' not found", extension)
}
