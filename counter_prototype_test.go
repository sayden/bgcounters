package counters

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPrototypes(t *testing.T) {
	// read the json file with testing data
	filepath := "./testdata/prototype.json"
	newTempl, err := ParseCountersJsonFile(filepath, "")
	if !assert.NoError(t, err, "could not parse the template") {
		t.FailNow()
	}

	// Extract the counters from the prototypes into a new template
	newTempl, err = ParsePrototypedTemplate(newTempl)
	assert.NoError(t, err)

	// check the new template
	assert.Equal(t, 6, len(newTempl.Counters))
	assert.Equal(t, 2, len(newTempl.Counters[0].Texts))

	// check the marshalling of the template to an expected byte slice
	actualBytes, err := json.MarshalIndent(newTempl, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, 4436, len(actualBytes))

	expectedFile, err := os.Open("./testdata/parse_template_01.json")
	assert.NoError(t, err)
	defer expectedFile.Close()

	expectedBytes, err := io.ReadAll(expectedFile)
	assert.NoError(t, err)
	// // ensure we are using the expected file and that it has not been altered by mistake
	if !assert.Equal(t, 4436, len(expectedBytes), "expected file has been altered, aborting test") {
		t.FailNow()
	}

	// compare the bytes of the expected file data and the actual data
	assert.Equal(t, len(expectedBytes), len(actualBytes), "Expected and actual byte slices have different lengths")
	assert.Equal(t, string(expectedBytes), string(actualBytes), "Expected and actual byte data are different")
}
