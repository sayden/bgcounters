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

func TestApplyCounterWaterfallSettings(t *testing.T) {
	ct := &CounterTemplate{
		Settings: Settings{
			Width:           100,
			Height:          200,
			FontHeight:      10,
			FontColorS:      "black",
			BackgroundColor: "white",
			Margins:         5,
			StrokeWidth:     1,
		},
		Counters: []Counter{
			{
				Settings: Settings{
					Margins:    10,
					FontHeight: 20,
				},
				Texts: []Text{{String: "text"}},
			},
		},
	}

	ct.ApplyCounterWaterfallSettings()

	assert.Equal(t, 20.0, ct.Counters[0].Settings.FontHeight)
	assert.Equal(t, "black", ct.Counters[0].Settings.FontColorS)
	assert.Equal(t, "white", ct.Counters[0].Settings.BackgroundColor)
	assert.Equal(t, 20.0, ct.Counters[0].Texts[0].Settings.FontHeight)
	assert.Equal(t, 100, ct.Counters[0].Settings.Width)
	assert.Equal(t, 200, ct.Counters[0].Settings.Height)
	assert.Equal(t, 10.0, ct.Counters[0].Settings.Margins)

	// Override with zero values
	ct.Counters[0].Settings.StrokeWidth = 0
	ct.Counters[0].Settings.Width = 50
	ct.ApplyCounterWaterfallSettings()

	assert.Equal(t, 0.0, ct.Counters[0].Settings.StrokeWidth, "StrokeWidth should be 1 for CT "+
		"settings but 0 for counter because it was overriden")
	assert.Equal(t, 50, ct.Counters[0].Settings.Width, "Width should be 50 for counter because it "+
		"was overriden")
}
