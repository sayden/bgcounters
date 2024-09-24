package fsops

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ReadMarkupFile(markupFilepath string, destination interface{}) error {
	extension := filepath.Ext(markupFilepath)

	data, err := ioutil.ReadFile(markupFilepath)
	if err != nil {
		return errors.Wrap(err, "could not read file content")
	}

	switch extension {
	case ".xml":
		if err = xml.Unmarshal(data, &destination); err != nil {
			return errors.Wrapf(err, "the file in '%s' has syntax errors", markupFilepath)
		}
		return nil
	case ".json":
		if err = json.Unmarshal(data, &destination); err != nil {
			return errors.Wrapf(err, "the file in '%s' has syntax errors", markupFilepath)
		}
		return nil
	}

	return fmt.Errorf("file extension '%s' not recognized. Use .json or .xml files only", extension)
}

func FilenameExistsInFolder(filename, folder string) bool {
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		log.WithError(err).Fatal("could not read images folder")
	}

	for _, file := range fs {
		_, existingFilename := filepath.Split(file.Name())
		_, gameImageName := filepath.Split(filename)
		if strings.Contains(gameImageName, existingFilename) {
			return true
		}
	}

	return false
}
