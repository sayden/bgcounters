package counters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaAtPath(t *testing.T) {
	err := ValidateSchemaAtPath[CounterTemplate]("./testdata/validate_schema_correct.json")
	assert.NoError(t, err)

	err = ValidateSchemaAtPath[CounterTemplate]("./testdata/validate_schema_incorrect.json")
	assert.Error(t, err)
}

func TestValidateSchemaReader(t *testing.T) {
	f, err := os.Open("./testdata/validate_schema_correct.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaReader[CounterTemplate](f)
		assert.NoError(t, err)
	}
	defer f.Close()

	f, err = os.Open("./testdata/validate_schema_incorrect.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaReader[CounterTemplate](f)
		assert.Error(t, err)
	}
	defer f.Close()

}

func TestValidateSchemaBytes(t *testing.T) {
	byt, err := os.ReadFile("./testdata/validate_schema_correct.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaBytes[CounterTemplate](byt)
		assert.NoError(t, err)
	}

	byt, err = os.ReadFile("./testdata/validate_schema_incorrect.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaBytes[CounterTemplate](byt)
		assert.Error(t, err)
	}
}
