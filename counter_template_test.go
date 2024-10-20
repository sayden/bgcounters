package counters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
