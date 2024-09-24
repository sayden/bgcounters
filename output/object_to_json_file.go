package output

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func ToJSONFile(input interface{}, outputPath string) error {
	dir := filepath.Dir(outputPath)
	_ = os.MkdirAll(dir, 0750)

	byt, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	genFile2, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "error trying to create file to write output")
	}
	defer genFile2.Close()

	_, err = genFile2.Write(byt)

	return err
}
