package transform

import (
	"io"
	"os"
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestJSONPrototypes(t *testing.T) {
	// read the json file with testing data
	filepath := "../testdata/prototype.json"
	f, err := os.Open(filepath)
	assert.NoError(t, err)
	defer f.Close()

	byt, err := io.ReadAll(f)
	assert.NoError(t, err)

	newTempl, err := counters.ParseTemplate(byt)
	assert.NoError(t, err)

	newTempl, err = JsonPrototypeToTemplate(newTempl)
	assert.NoError(t, err)

	// check the new template
	assert.Equal(t, 2, len(newTempl.Counters))
	assert.Equal(t, 1, len(newTempl.Counters[0].Texts))
}
