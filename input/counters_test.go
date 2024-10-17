package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCounterTemplate(t *testing.T) {
	inputPath1 := "../testdata/validate_schema_correct.json"
	inputPathCsv := "../testdata/testing_cards.csv"

	_, err := ReadCounterTemplate(inputPath1)
	assert.NoError(t, err)

	_, err = ReadCounterTemplate(inputPathCsv)
	assert.NoError(t, err)
}
