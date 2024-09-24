package fsops

import (
	"path/filepath"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

func GetExtension(cardsFilepath string) (counters.FileContent, error) {
	extension := filepath.Ext(cardsFilepath)

	switch extension {
	case ".csv":
		return counters.FileContent_CSV, nil
	case ".json":
		return counters.FileContent_JSON, nil
	}

	return 0, errors.Errorf("file extension '%s' not recognized", extension)
}

func IdentifyJSONFileContent(data []byte) (counters.FileContent, error) {
	_, _, _, err := jsonparser.Get(data, "cards")
	if err == nil {
		return counters.FileContent_CardTemplate, nil
	}

	_, _, _, err = jsonparser.Get(data, "counters")
	if err == nil {
		return counters.FileContent_CounterTemplate, nil
	}

	_, dataType, _, err := jsonparser.Get(data)
	if err != nil {
		return 0, errors.Wrap(err, "could not get root data in file")
	}

	var dataTypeInt counters.FileContent = -1
	switch dataType {
	case jsonparser.Array:
		// It can be an events file or a quotes file
		_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, _, _, err = jsonparser.Get(value, "title")
			if err == nil {
				dataTypeInt = counters.FileContent_Events
			}

			_, _, _, err = jsonparser.Get(value, "quote")
			if err == nil {
				dataTypeInt = counters.FileContent_Quotes
			}
		})

		if err != nil {
			return 0, err
		}
	}

	if dataTypeInt != -1 {
		return dataTypeInt, nil
	}

	return 0, errors.New("the content of the file is not recognized as card, counters, quotes or events")
}
