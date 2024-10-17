package counters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaAtPath(t *testing.T) {
	err := ValidateSchemaAtPath("./testdata/validate_schema_correct.json")
	assert.NoError(t, err)

	err = ValidateSchemaAtPath("./testdata/validate_schema_incorrect.json")
	assert.Error(t, err)
}

func TestValidateSchemaReader(t *testing.T) {
	f, err := os.Open("./testdata/validate_schema_correct.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaReader(f)
		assert.NoError(t, err)
	}
	defer f.Close()

	f, err = os.Open("./testdata/validate_schema_incorrect.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaReader(f)
		assert.Error(t, err)
	}
	defer f.Close()

}

func TestValidateSchemaBytes(t *testing.T) {
	byt, err := os.ReadFile("./testdata/validate_schema_correct.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaBytes(byt)
		assert.NoError(t, err)
	}

	byt, err = os.ReadFile("./testdata/validate_schema_incorrect.json")
	if assert.NoError(t, err) {
		err = ValidateSchemaBytes(byt)
		assert.Error(t, err)
	}
}

func TestExpandPrototypeCounterTemplate(t *testing.T) {
	proto := CounterPrototype{
		Counter: Counter{
			Texts: []Text{{String: "text"}},
		},
		TextPrototypes: []TextPrototype{
			{StringList: []string{"text1", "text2"}},
		},
		ImagePrototypes: []ImagePrototype{
			{PathList: []string{"../assets/binoculars.png", "../assets/stripe.png"}},
		},
	}

	prototypeTemplate := &CounterTemplate{
		Prototypes: map[string]CounterPrototype{
			"proto":  proto,
			"proto2": proto,
		}}

	ct, err := prototypeTemplate.ExpandPrototypeCounterTemplate()
	if assert.NoError(t, err) {
		assert.Equal(t, 4, len(ct.Counters))
		assert.Equal(t, "text", ct.Counters[0].Texts[0].String)
		assert.Equal(t, "text1", ct.Counters[0].Texts[1].String)
		assert.Equal(t, "../assets/binoculars.png", ct.Counters[0].Images[0].Path)

		assert.Equal(t, "text", ct.Counters[1].Texts[0].String)
		assert.Equal(t, "text2", ct.Counters[1].Texts[1].String)
		assert.Equal(t, "../assets/stripe.png", ct.Counters[1].Images[0].Path)

		assert.Equal(t, "text", ct.Counters[2].Texts[0].String)
		assert.Equal(t, "text1", ct.Counters[2].Texts[1].String)
		assert.Equal(t, "../assets/binoculars.png", ct.Counters[2].Images[0].Path)

		assert.Equal(t, "text", ct.Counters[3].Texts[0].String)
		assert.Equal(t, "text2", ct.Counters[3].Texts[1].String)
		assert.Equal(t, "../assets/stripe.png", ct.Counters[3].Images[0].Path)
	}
}
